package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

func (r *TodoItemPostgres) Create(listId int, todoItem todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}

	var itemId int
	// We have insert user's data to todo_items table
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	// method QueryRow return Row structure that has row or error.
	row := tx.QueryRow(createItemQuery, todoItem.Title, todoItem.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// We have to put data as well in table lists_items item_id that correspondent to list_id
	createListsItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListsItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(id, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.* FROM %s ti INNER JOIN %s li ON ti.id = li.item_id
							INNER JOIN %s ul ON li.list_id = ul.list_id 
							WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, id); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(id, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	query := fmt.Sprintf(`SELECT ti.* FROM %s ti INNER JOIN %s li ON ti.id = li.item_id
						 INNER JOIN %s ul ON li.list_id = ul.list_id 
						 WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&item, query, itemId, id); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemPostgres) Delete(id, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.list_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2",
		todoItemsTable, listsItemsTable, usersListsTable)
	result, err := r.db.Exec(query, id, itemId)
	if int, err := result.RowsAffected(); err !=nil {
		if int == 0 {
			return errors.New("Nothing to delete")
		}
	}

	return err
}

func (r *TodoItemPostgres) Update(id, itemId int, input todo.UpdateItemInput) error {
		// init empty slice
		setValues := make([]string, 0)
		// init empty interface slice
		args := make([]interface{}, 0)
		argId := 1
	
		if input.Title != nil {
			setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
			args = append(args, *input.Title)
			argId++
		}
	
		if input.Description != nil {
			setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
			args = append(args, *input.Description)
			argId++
		}

		if input.Done != nil {
			setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
			args = append(args, *input.Done)
			argId++
		}
	
		//"title=$1, description=$2, done=$3"
		setQuery := strings.Join(setValues, ", ")
	
		query := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d",
			todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

		//[cook, some receipe, true, 7, 1]   exemple
		args = append(args, itemId, id)
	
		_, err := r.db.Exec(query, args...)
		return err
}

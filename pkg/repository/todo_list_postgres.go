package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MaximumTroubles/go-todo-app"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(id int, todoList todo.TodoList) (int, error) {
	// we have to do here 2 operations, insert todo list to table todo_list, and insert list and user id in table users_list
	// Here comes transaction!

	// To start transaction method use method Begin()
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// first we need to save todo list and return Id created list
	var listId int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, todoList.Title, todoList.Description)
	if err := row.Scan(&listId); err != nil {
		// if we can't find id by method Scan it mean something wrong and err happend, so we provide Rollback()
		tx.Rollback()
		return 0, err
	}

	// second we need to save user_id and list_id in users_lists
	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	// for simple query without reading returning info we can use method Exec()
	_, err = tx.Exec(createUsersListsQuery, id, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// when everything went well call Commit()
	tx.Commit()
	return listId, nil
}

func (r *TodoListPostgres) GetAll(id int) ([]todo.TodoList, error) {
	// init todo.TodoList entity slice.
	var lists []todo.TodoList

	query := fmt.Sprintf("SELECT tl.* FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	// Here we use method Select because it pecieves slices, and select more than one item from db.
	err := r.db.Select(&lists, query, id)

	return lists, err
}

func (r *TodoListPostgres) GetById(id, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("SELECT tl.* FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, id, listId)

	return list, err
}

func (r *TodoListPostgres) Update(id, listId int, input todo.UpdateListInput) error {
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

	//"title=$1, description=$2"
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, id)

	fmt.Printf("updateQuery: %s\n", query)
	fmt.Printf("args: %s\n", args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoListPostgres) Delete(id, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	result, err := r.db.Exec(query, id, listId)
	rowAffected, err := result.RowsAffected()
	if rowAffected == 0 {
		return errors.New("nothing to delete")
	}

	return err
}

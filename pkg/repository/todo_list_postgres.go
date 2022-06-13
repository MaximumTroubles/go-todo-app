package repository

import (
	"fmt"

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


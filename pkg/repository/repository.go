package repository

import (
	"github.com/MaximumTroubles/go-todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(id int, todoList todo.TodoList) (int, error)
	GetAll(id int) ([]todo.TodoList, error)
	GetById(id, listId int) (todo.TodoList, error)

}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList: NewTodoListPostgres(db),
	}
}

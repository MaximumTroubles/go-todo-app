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
	Delete(id,listId int) error
	Update(id, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, todoItem todo.TodoItem) (int, error)
	GetAll(id, listId int) ([]todo.TodoItem, error)
	GetById(id, itemId int) (todo.TodoItem, error)
	Delete(id, itemId int) error
	Update(id, itemId int, input todo.UpdateItemInput) error
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
		TodoItem: NewTodoItemPostgres(db),
	}
}

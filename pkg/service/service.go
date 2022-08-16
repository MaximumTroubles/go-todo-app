package service

import (
	"github.com/MaximumTroubles/go-todo-app"
	"github.com/MaximumTroubles/go-todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(id int, todoList todo.TodoList) (int, error)
	// when you need to return more than one item init slice of structurs
	GetAll(id int) ([]todo.TodoList, error)
	GetById(id, listId int) (todo.TodoList, error)
	Delete(id, listId int) error
	Update(id, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
	Create(id, listId int, todoItem todo.TodoItem) (int, error)
	GetAll(id, listId int) ([]todo.TodoItem, error)
	GetById(id, itemId int) (todo.TodoItem, error)
	Delete(id, itemId int) error
	Update(id, itemId int, input todo.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}

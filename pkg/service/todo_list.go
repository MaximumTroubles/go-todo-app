package service

import (
	"github.com/MaximumTroubles/go-todo-app"
	"github.com/MaximumTroubles/go-todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(id int, todoList todo.TodoList) (int, error) {
	return s.repo.Create(id, todoList)
}

func (s *TodoListService) GetAll(id int) ([]todo.TodoList, error) {
	return s.repo.GetAll(id)
}

func (s *TodoListService) GetById(id, listId int) (todo.TodoList, error) {
	return s.repo.GetById(id, listId)	
}

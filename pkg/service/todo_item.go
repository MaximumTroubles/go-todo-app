package service

import (
	"github.com/MaximumTroubles/go-todo-app"
	"github.com/MaximumTroubles/go-todo-app/pkg/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo: repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(id, listId int, todoItem todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(id, listId)
	if err != nil {
		// list does not exists or does not belongs to user
		return 0, err
	}

	return s.repo.Create(listId, todoItem)
}

func (s *TodoItemService) GetAll(id, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(id, listId)
}

func (s *TodoItemService) GetById(id, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(id, itemId)
}

func(s *TodoItemService) Delete(id, itemId int) error {
	return s.repo.Delete(id, itemId)
}

func (s *TodoItemService) Update(id, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(id, itemId, input)
}

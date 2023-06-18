package service

import (
	"toDoList"
	"toDoList/pkg/repository"
)

type TodoItemService struct {
	repo     repository.ToDoItem
	listRepo repository.ToDoList
}

func NewTodoItemService(repo repository.ToDoItem, listRepo repository.ToDoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item toDoList.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId int, listId int) ([]toDoList.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId int, itemId int) (toDoList.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId int, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId int, itemId int, input toDoList.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, itemId, input)
}

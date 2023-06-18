package service

import (
	"toDoList"
	"toDoList/pkg/repository"
)

type Authorization interface {
	CreateUser(user toDoList.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type ToDoList interface {
	Create(userId int, list toDoList.TodoList) (int, error)
	GetAll(userId int) ([]toDoList.TodoList, error)
	GetById(userId int, listId int) (toDoList.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, input toDoList.UpdateListInput) error
}

type ToDoItem interface {
	Create(userId int, listId int, item toDoList.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]toDoList.TodoItem, error)
	GetById(userId int, itemId int) (toDoList.TodoItem, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input toDoList.UpdateItemInput) error
}

type Service struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: newAuthService(repos.Authorization),
		ToDoList: NewTodoListService(repos.ToDoList),
		ToDoItem: NewTodoItemService(repos.ToDoItem, repos.ToDoList)}
}

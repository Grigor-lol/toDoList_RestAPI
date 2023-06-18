package repository

import (
	"github.com/jmoiron/sqlx"
	"toDoList"
)

type Authorization interface {
	CreateUser(user toDoList.User) (int, error)
	GetUser(username, password string) (toDoList.User, error)
}

type ToDoList interface {
	Create(userId int, list toDoList.TodoList) (int, error)
	GetAll(userId int) ([]toDoList.TodoList, error)
	GetById(userId int, listId int) (toDoList.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, input toDoList.UpdateListInput) error
}

type ToDoItem interface {
	Create(listId int, item toDoList.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]toDoList.TodoItem, error)
	GetById(userId int, itemId int) (toDoList.TodoItem, error)
	Delete(userId int, itemId int) error
	Update(userId int, itemId int, input toDoList.UpdateItemInput) error
}

type Repository struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthMysql(db),
		ToDoList:      NewTodoListMySql(db),
		ToDoItem:      NewTodoItemMySql(db)}
}

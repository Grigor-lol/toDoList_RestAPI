package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
	"toDoList"
)

type TodoListMySql struct {
	db *sqlx.DB
}

func NewTodoListMySql(db *sqlx.DB) *TodoListMySql {
	return &TodoListMySql{db: db}
}

func (r *TodoListMySql) Create(userId int, list toDoList.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES (?, ?)", todoListsTable)
	row, err := tx.Exec(createListQuery, list.Title, list.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var id int64

	id, err = row.LastInsertId()
	if err != nil {
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES (?, ?)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), tx.Commit()
}

func (r *TodoListMySql) GetAll(userId int) ([]toDoList.TodoList, error) {
	var lists []toDoList.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = ?",
		todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListMySql) GetById(userId int, listId int) (toDoList.TodoList, error) {
	var list toDoList.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = ? AND ul.list_id = ?",
		todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListMySql) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE tl FROM %s tl JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = ? AND ul.list_id = ?",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListMySql) Update(userId, listId int, input toDoList.UpdateListInput) error {
	setValues := make([]string, 0)

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=\"%s\"", *input.Title))
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description= \"%s\"", *input.Description))
	}

	setQuery := strings.Join(setValues, ", ")
	logrus.Debugf(setQuery)

	query := fmt.Sprintf("UPDATE %s tl JOIN %s ul ON tl.id = ul.list_id SET %s WHERE ul.list_id=? AND ul.user_id=?",
		todoListsTable, usersListsTable, setQuery)

	_, err := r.db.Exec(query, listId, userId)
	return err
}

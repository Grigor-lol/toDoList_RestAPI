package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"toDoList"
)

type TodoItemMySql struct {
	db *sqlx.DB
}

func NewTodoItemMySql(db *sqlx.DB) *TodoItemMySql {
	return &TodoItemMySql{db: db}
}

func (r *TodoItemMySql) Create(listId int, item toDoList.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int64
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values (?, ?)", todoItemsTable)

	row, err := tx.Exec(createItemQuery, item.Title, item.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	itemId, err = row.LastInsertId()
	if err != nil {
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values (?, ?)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(itemId), tx.Commit()
}
func (r *TodoItemMySql) GetAll(userId, listId int) ([]toDoList.TodoItem, error) {
	var items []toDoList.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = ? AND ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemMySql) GetById(userId int, itemId int) (toDoList.TodoItem, error) {
	var item toDoList.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = ? AND ul.user_id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil

}
func (r *TodoItemMySql) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE ti FROM %s ti JOIN %s li ON ti.id = li.item_id
													  JOIN %s ul ON li.list_id = ul.list_id
									WHERE ul.user_id = ? AND ti.id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, itemId)
	return err
}

func (r *TodoItemMySql) Update(userId, itemId int, input toDoList.UpdateItemInput) error {
	setValues := make([]string, 0)

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=\"%s\"", *input.Title))
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description= \"%s\"", *input.Description))
	}

	var done = 0
	if input.Done != nil {
		if *input.Done {
			done = 1
		}
		setValues = append(setValues, fmt.Sprintf("done= %d", done))
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti 
								JOIN %s li ON ti.id = li.item_id
								JOIN %s ul ON li.list_id = ul.list_id
								SET %s 
								WHERE ul.user_id = ? AND ti.id = ?`,
		todoItemsTable, listsItemsTable, usersListsTable, setQuery)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

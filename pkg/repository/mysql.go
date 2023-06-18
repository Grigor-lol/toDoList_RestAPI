package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Config struct {
	Username string
	Password string
}

func NewMySqlDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", cfg.Username+":"+cfg.Password+"@tcp(localhost:3306)/todolist")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

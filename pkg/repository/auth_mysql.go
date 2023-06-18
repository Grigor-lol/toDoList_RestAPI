package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"toDoList"
)

type AuthMysql struct {
	db *sqlx.DB
}

func newAuthMysql(db *sqlx.DB) *AuthMysql {
	return &AuthMysql{
		db: db,
	}
}

func (r *AuthMysql) CreateUser(user toDoList.User) (int, error) {
	//var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values (?, ?, ?)", usersTable)

	r.db.QueryRow(query, user.Name, user.Username, user.Password)
	//if err := row.Scan(&id); err != nil {
	//	return 0, err
	//}
	return 0, nil
}

func (r *AuthMysql) GetUser(username, password string) (toDoList.User, error) {
	var user toDoList.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=? AND password_hash=?", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

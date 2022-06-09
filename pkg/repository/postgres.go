package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	todoListstsTable = "todo_lists"
	usersListsTable  = "users_lists"
	todoItemsTable   = "todo_items"
	listsItemsTable  = "lists_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBname, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}
	// as far as I saw inside Connect method, it also make Ping function so it's not necessary to do it here
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

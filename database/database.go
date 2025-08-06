package database

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		return nil, err
	}

	if err := CreateUsersTable(db); err != nil {
		return nil, err
	}

	if err := CreateMessagesTable(db); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`)
	return err
}

func CreateMessagesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			message TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)
	return err
}

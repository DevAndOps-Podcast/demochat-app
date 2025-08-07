package database

import (
	"context"
	"database/sql"
	"demochat/config"
	"fmt"

	_ "github.com/lib/pq"
)

func New(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Set the schema for the current session
	_, err = db.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s", cfg.DB.Schema))
	if err != nil {
		return nil, err
	}

	if err := CreateSchema(db, cfg.DB.Schema); err != nil {
		return nil, err
	}

	if err := CreateUsersTable(db, cfg.DB.Schema); err != nil {
		return nil, err
	}

	if err := CreateMessagesTable(db, cfg.DB.Schema); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateSchema(db *sql.DB, schemaName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	return err
}

func CreateUsersTable(db *sql.DB, schemaName string) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`, schemaName))
	return err
}

func CreateMessagesTable(db *sql.DB, schemaName string) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.messages (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			message TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES %s.users(id)
		);
	`, schemaName, schemaName))
	return err
}

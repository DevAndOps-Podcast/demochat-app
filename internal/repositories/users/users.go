package users

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64
	Username string
	Password string
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	row := r.db.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE username = ?", username)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User

	row := r.db.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE id = ?", id)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	return err
}

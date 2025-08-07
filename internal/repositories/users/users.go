package users

import (
	"context"
	"database/sql"
	"demochat/config"
	"fmt"
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
	db     *sql.DB
	schema string
}

func New(db *sql.DB, cfg *config.Config) Repository {
	return &repository{db: db, schema: cfg.DB.Schema}
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT id, username, password FROM %s.users WHERE username = $1", r.schema), username)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User

	row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT id, username, password FROM %s.users WHERE id = $1", r.schema), id)

	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s.users (username, password) VALUES ($1, $2)", r.schema), user.Username, user.Password)
	return err
}

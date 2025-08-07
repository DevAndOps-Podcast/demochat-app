package messages

import (
	"context"
	"database/sql"
	"demochat/config"
	"fmt"
)

type Message struct {
	ID       int64
	UserID   int64
	Username string
	Message  string
}

type Repository interface {
	CreateMessage(ctx context.Context, message *Message) error
	ListMessages(ctx context.Context) ([]*Message, error)
}

type repository struct {
	db     *sql.DB
	schema string
}

func New(db *sql.DB, cfg *config.Config) Repository {
	return &repository{db: db, schema: cfg.DB.Schema}
}

func (r *repository) CreateMessage(ctx context.Context, message *Message) error {
	_, err := r.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s.messages (user_id, message) VALUES ($1, $2)", r.schema), message.UserID, message.Message)
	return err
}

func (r *repository) ListMessages(ctx context.Context) ([]*Message, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT m.id, m.user_id, u.username, m.message FROM %s.messages m INNER JOIN %s.users u ON m.user_id = u.id ORDER BY m.id ASC", r.schema, r.schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.UserID, &msg.Username, &msg.Message); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

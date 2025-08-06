package messages

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
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
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMessage(ctx context.Context, message *Message) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO messages (user_id, message) VALUES (?, ?)", message.UserID, message.Message)
	return err
}

func (r *repository) ListMessages(ctx context.Context) ([]*Message, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT m.id, m.user_id, u.username, m.message FROM messages m INNER JOIN users u ON m.user_id = u.id ORDER BY m.id ASC")
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

package repositories

import (
	"database/sql"
	"demochat/config"
	"demochat/internal/repositories/messages"
	"demochat/internal/repositories/users"

	"go.uber.org/fx"
)

var Module = fx.Module("repositories", fx.Provide(
	func(db *sql.DB, cfg *config.Config) users.Repository {
		return users.New(db, cfg)
	},
	func(db *sql.DB, cfg *config.Config) messages.Repository {
		return messages.New(db, cfg)
	},
))

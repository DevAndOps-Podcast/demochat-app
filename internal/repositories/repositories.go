package repositories

import (
	"demochat/internal/repositories/messages"
	"demochat/internal/repositories/users"

	"go.uber.org/fx"
)

var Module = fx.Module("repositories", fx.Provide(
	users.New,
	messages.New,
))

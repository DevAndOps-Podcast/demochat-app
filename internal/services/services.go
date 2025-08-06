package services

import (
	"demochat/internal/services/auth"
	"demochat/internal/services/insights"
	"demochat/internal/services/messages"

	"go.uber.org/fx"
)

var Module = fx.Module("services", fx.Provide(
	auth.New,
	messages.New,
	insights.New,
))

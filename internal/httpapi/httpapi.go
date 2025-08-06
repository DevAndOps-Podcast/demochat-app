package httpapi

import (
	"demochat/internal/httpapi/handlers/auth"
	"demochat/internal/httpapi/handlers/healthcheck"
	"demochat/internal/httpapi/handlers/insights"
	"demochat/internal/httpapi/handlers/messages"
	"demochat/internal/httpapi/handlers/static"

	"go.uber.org/fx"
)

var Module = fx.Module("httpapi", fx.Provide(
	auth.New,
	static.New,
	messages.New,
	healthcheck.New,
	insights.New,
))

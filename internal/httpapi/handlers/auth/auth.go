package auth

import (
	"demochat/internal/httpapi/handlers"
	"demochat/internal/services/auth"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct {
	AuthService auth.Service
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

func New(s auth.Service) Result {
	return Result{
		Handler: &handler{
			AuthService: s,
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/auth")
	group.POST("", h.Authenticate)
	group.POST("/register", h.Register)
	group.POST("/token/refresh", h.RefreshToken)
}

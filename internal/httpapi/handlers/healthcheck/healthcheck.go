package healthcheck

import (
	"demochat/internal/httpapi/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct{}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

func New() Result {
	return Result{
		Handler: &handler{},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)
}

func (h *handler) HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

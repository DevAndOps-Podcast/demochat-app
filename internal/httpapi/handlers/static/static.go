package static

import (
	"demochat/internal/httpapi/handlers"

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
	e.Static("/static", "static")
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})
}

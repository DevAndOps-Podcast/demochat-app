package static

import (
	"demochat/internal/httpapi/handlers"
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

//go:embed content
var content embed.FS

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
	// Create a subdirectory for the embedded content
	staticContent, _ := fs.Sub(content, "static")
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(staticContent)))))
}

package main

import (
	"context"
	"database/sql"
	"demochat/clients/insights"
	"demochat/config"
	"demochat/database"
	"demochat/internal/httpapi"
	"demochat/internal/httpapi/handlers"
	"demochat/internal/repositories"
	"demochat/internal/services"
	"demochat/logger"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Lc       fx.Lifecycle
	Config   *config.Config
	DB       *sql.DB
	Echo     *echo.Echo
	Handlers []handlers.Handler `group:"handlers"`
	Logger   *logger.Logger
}

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			config.New,
			logger.New,
			database.New,
			echo.New,
			insights.New,
		),
		repositories.Module,
		services.Module,
		httpapi.Module,
		fx.Invoke(
			setLifeCycle,
		),
	)

	app.Run()
}

func setLifeCycle(p Params) {
	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			p.Echo.Use(middleware.Recover())
			p.Echo.Use(middleware.RequestID())
			p.Echo.Use(middleware.Logger())

			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo)
			}

			go func() {
				p.Echo.Logger.Error(p.Echo.Start(p.Config.Address))
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := p.Echo.Shutdown(ctx); err != nil {
				p.Logger.Error("failed to shutdown echo", slog.Any("error", err))
			}

			if err := p.DB.Close(); err != nil {
				p.Logger.Error("failed to close database connection", slog.Any("error", err))
			}

			return nil
		},
	})
}

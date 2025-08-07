package insights

import (
	"demochat/internal/httpapi/handlers"
	"demochat/internal/services/auth"
	"demochat/internal/services/insights"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct {
	insightsService insights.Service
	authService     auth.Service
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

func New(svc insights.Service, authService auth.Service) Result {
	return Result{
		Handler: &handler{
			insightsService: svc,
			authService:     authService,
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/insights", h.GetInsights)
}

type InsightsResponse struct {
	TotalMessages      int64   `json:"total_messages"`
	MostActiveUser     string  `json:"most_active_user"`
	AverageMessageRate float64 `json:"average_message_rate"`
}

func (h *handler) GetInsights(c echo.Context) error {
	insights := h.insightsService.GetInsights(c.Request().Context())
	username := ""

	log.Println(insights)

	u, err := h.authService.FindByID(c.Request().Context(), insights.MostActiveUserID)
	if err != nil {
		log.Println("failed to load most_active_user username with id:", insights.MostActiveUserID)
	} else if u != nil {
		username = u.Username
	}

	resp := InsightsResponse{
		TotalMessages:      insights.TotalMessages,
		MostActiveUser:     username,
		AverageMessageRate: insights.AverageMessageRate,
	}

	return c.JSON(http.StatusOK, resp)
}

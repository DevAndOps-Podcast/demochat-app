package messages

import (
	"demochat/config"
	"demochat/internal/httpapi/handlers"
	"demochat/internal/services/messages"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type handler struct {
	MessagesService messages.Service
	JWTSecret       string
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type CreateMessageRequest struct {
	Message string `json:"message"`
}

func New(messagesService messages.Service, cfg *config.Config) Result {
	return Result{
		Handler: &handler{
			MessagesService: messagesService,
			JWTSecret:       cfg.JWTSecret,
		},
	}
}

func (h *handler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/messages")
	group.Use(echojwt.JWT([]byte(h.JWTSecret)))
	group.POST("", h.CreateMessage)
	group.GET("", h.ListMessages)
}

func (h *handler) CreateMessage(c echo.Context) error {
	var req CreateMessageRequest

	if err := c.Bind(&req); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["sub"].(float64))

	if err := h.MessagesService.SaveMessage(c.Request().Context(), userID, req.Message); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

type MessageResponse struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func (h *handler) ListMessages(c echo.Context) error {
	messages, err := h.MessagesService.ListMessages(c.Request().Context())
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response []MessageResponse
	for _, msg := range messages {
		response = append(response, MessageResponse{
			ID:       msg.ID,
			UserID:   msg.UserID,
			Username: msg.Username,
			Message:  msg.Message,
		})
	}

	return c.JSON(http.StatusOK, response)
}

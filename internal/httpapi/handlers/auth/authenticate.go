package auth

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Authenticate(c echo.Context) error {
	var req AuthenticationRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log.Println(req.Username, req.Password)
	res, err := h.AuthService.Authenticate(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, &AuthenticationResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

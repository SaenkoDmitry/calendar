package controllers

import (
	"calendar/internal/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) HealthCheck(c echo.Context) error {

	// TODO check services

	return c.JSON(http.StatusOK, models.HealthStatus{
		Status:    "OK",
		Timestamp: time.Now(),
	})
}

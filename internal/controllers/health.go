package controllers

import (
	"calendar/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *handler) HealthCheck(c echo.Context) error {
	dbStatus := "DOWN"
	if h.pool != nil {
		if err := h.pool.Ping(context.Background()); err == nil {
			dbStatus = "UP"
		}
	}

	return c.JSON(http.StatusOK, models.HealthStatus{
		Status:    "OK",
		Timestamp: time.Now(),
		DB:        dbStatus,
	})
}

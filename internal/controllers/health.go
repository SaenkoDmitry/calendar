package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary  service health info
// @Tags     health
// @Accept   json
// @Produce  json
// @Success  200  {object}  models.HealthStatus
// @Router   /health [get]
func (h *handler) HealthCheck(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, models.HealthStatus{
		Status:    "OK",
		Timestamp: time.Now(),
		Services: []*models.HealthService{
			{
				Name:   constants.PostgresDBService,
				Status: h.DB.Check(context.Background()),
			},
		},
	}, "  ")
}

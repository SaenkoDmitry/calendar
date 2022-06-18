package controllers

import (
	"calendar/internal/helpers"
	"calendar/internal/models"
	"calendar/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) ChangeUserZone(c echo.Context) error {
	req := new(models.ChangeZoneReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	userID, err := helpers.GetUser(c, "userID")
	if err != nil {
		return err
	}

	loc, err := helpers.ChooseZone(c, req.Zone)
	if err != nil {
		return err
	}

	err = repository.UpdateUserZone(c, h.pool, userID, loc)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusOK, loc.String())
}

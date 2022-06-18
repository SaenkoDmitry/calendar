package controllers

import (
	"calendar/internal/helpers"
	"calendar/internal/models"
	"calendar/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) CreateUser(c echo.Context) error {
	req := new(models.CreateUserReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	loc, err := helpers.ChooseZone(c, req.Zone)
	if err != nil {
		return err
	}

	err = repository.CreateUser(c, h.pool, req.FirstName, req.SecondName, req.Email, loc)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusCreated, req.Email)
}

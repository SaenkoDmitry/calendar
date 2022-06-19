package controllers

import (
	"calendar/internal/helpers"
	"calendar/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateUser godoc
// @Summary  create new user
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    CreateUserReq  body      models.CreateUserReq  true  "request body for creating user"
// @Success  201            {object}  models.DataError
// @Failure  400            {object}  models.DataError
// @Failure  500            {object}  models.DataError
// @Router   /users [post]
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

	_, err = h.DB.CreateUser(c, req.FirstName, req.SecondName, req.Email, loc)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusCreated, req.Email)
}

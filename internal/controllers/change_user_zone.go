package controllers

import (
	"calendar/internal/helpers"
	"calendar/internal/models"
	"calendar/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ChangeUserZone godoc
// @Summary  change user time zone
// @Tags     user
// @Accept   json
// @Produce  json
// @Param    CreateUserReq  body      models.ChangeZoneReq  true  "request body for changing user time zone"
// @Success  200            {object}  models.DataError
// @Failure  400            {object}  models.DataError
// @Failure  500            {object}  models.DataError
// @Router   /users/{userID} [put]
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

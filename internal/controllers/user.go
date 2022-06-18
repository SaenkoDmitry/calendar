package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"net/http"
	"time"

	"github.com/jackc/pgconn"

	"github.com/labstack/echo/v4"
)

func (h *handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
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

	_, err = h.pool.Exec(ctx, "INSERT INTO users(first_name, second_name, email) "+
		"	VALUES($1, $2, $3)", req.FirstName, req.SecondName, req.Email, loc.String())
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if helpers.IsEmailDuplicated(pgErr) {
			return helpers.WrapError(c, http.StatusConflict, constants.EmailAlreadyRegistered)
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return helpers.WrapSuccess(c, http.StatusCreated, req.Email)
}

func (h *handler) ChangeUserZone(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.ChangeZoneReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	userID := c.Param("id")

	loc, err := time.LoadLocation(req.Zone)
	if err != nil {
		return helpers.WrapError(c, http.StatusBadRequest, constants.NotValidTimeZone)
	}

	res, err := h.pool.Exec(ctx, "UPDATE users SET user_zone = $1 WHERE id = $2", loc.String(), userID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if res.RowsAffected() == 0 {
		return helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	_ = res

	return helpers.WrapSuccess(c, http.StatusOK, loc)
}

func (h *handler) GetMeetingsByUser(c echo.Context) error {

	return nil
}

func (h *handler) GetMeetingsForGroupOfUsers(c echo.Context) error {

	return nil
}

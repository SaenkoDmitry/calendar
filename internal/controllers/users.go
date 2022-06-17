package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgconn"

	"github.com/labstack/echo/v4"
)

type CreateUserReq struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
}

type ChangeZoneReq struct {
	Zone string `json:"zone" validate:"required"`
}

func (h *handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(CreateUserReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	_, err := h.pool.Exec(ctx, "INSERT INTO users(first_name, second_name, email) "+
		"	VALUES($1, $2, $3)", req.FirstName, req.SecondName, req.Email)
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if isEmailDuplicated(pgErr) {
			return utils.WrapJSONError(c, http.StatusConflict, constants.EmailAlreadyRegistered)
		}
		return utils.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return c.JSON(http.StatusCreated, req.Email)
}

func (h *handler) ChangeUserZone(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(ChangeZoneReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}
	userID := c.Param("id")

	loc, err := time.LoadLocation(req.Zone)
	if err != nil {
		return utils.WrapJSONError(c, http.StatusBadRequest, constants.NotValidTimeZone)
	}

	_, err = h.pool.Exec(ctx, "UPDATE users SET user_zone = $1 WHERE id = $2", loc.String(), userID)
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if true {
			_ = pgErr
			return utils.WrapJSONError(c, http.StatusBadRequest, constants.UserIDNotExists)
		}
		return utils.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return c.JSON(http.StatusOK, loc)
}

func isEmailDuplicated(err *pgconn.PgError) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Message, "duplicate key")
}

func (h *handler) GetMeetingsByUser(c echo.Context) error {

	return nil
}

func (h *handler) GetMeetingsForGroupOfUsers(c echo.Context) error {

	return nil
}

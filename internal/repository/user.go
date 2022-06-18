package repository

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"net/http"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	selectUserZoneSQL = `SELECT user_zone FROM users WHERE id = $1`
	updateUserZoneSQL = `UPDATE users SET user_zone = $1 WHERE id = $2`
	createUserSQL     = `INSERT INTO users(first_name, second_name, email, user_zone) VALUES($1, $2, $3, $4)`
)

func SelectUserZone(c echo.Context, pool *pgxpool.Pool, userID int32) (*time.Location, error) {
	var zoneStr string
	ctx := c.Request().Context()
	row := pool.QueryRow(ctx, selectUserZoneSQL, userID)
	if err := row.Scan(&zoneStr); err != nil {
		return nil, helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}
	loc, err := helpers.ChooseZone(c, zoneStr)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func UpdateUserZone(c echo.Context, pool *pgxpool.Pool, userID int32, loc *time.Location) error {
	ctx := c.Request().Context()
	res, err := pool.Exec(ctx, updateUserZoneSQL, loc.String(), userID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if res.RowsAffected() == 0 {
		return helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	return nil
}

func CreateUser(c echo.Context, pool *pgxpool.Pool, firstName, secondName, email string, loc *time.Location) error {
	ctx := c.Request().Context()
	_, err := pool.Exec(ctx, createUserSQL, firstName, secondName, email, loc.String())
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if helpers.IsEmailDuplicated(pgErr) {
			err = helpers.WrapError(c, http.StatusConflict, constants.EmailAlreadyRegistered)
			return err
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	return nil
}

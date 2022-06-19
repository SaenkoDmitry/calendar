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
	createUserSQL     = `INSERT INTO users(first_name, second_name, email, user_zone) VALUES($1, $2, $3, $4) RETURNING id`
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

func CreateUser(c echo.Context, pool *pgxpool.Pool, firstName, secondName, email string, loc *time.Location) (userID int32, err error) {
	ctx := c.Request().Context()
	row := pool.QueryRow(ctx, createUserSQL, firstName, secondName, email, loc.String())
	err = row.Scan(&userID)
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if helpers.IsEmailDuplicated(pgErr) {
			return 0, helpers.WrapError(c, http.StatusConflict, constants.EmailAlreadyRegistered)
		}
		return 0, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	return 0, nil
}

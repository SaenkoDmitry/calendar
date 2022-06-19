package repository

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"net/http"
	"time"

	"github.com/jackc/pgconn"

	"github.com/labstack/echo/v4"
)

const (
	selectUserZoneSQL = `SELECT user_zone FROM users WHERE id = $1`
	updateUserZoneSQL = `UPDATE users SET user_zone = $1 WHERE id = $2`
	createUserSQL     = `INSERT INTO users(first_name, second_name, email, user_zone) VALUES($1, $2, $3, $4) RETURNING id`
)

func (db *DB) SelectUserZone(c echo.Context, userID int32) (*time.Location, error) {
	var zoneStr string
	ctx := c.Request().Context()
	row := db.pool.QueryRow(ctx, selectUserZoneSQL, userID)
	if err := row.Scan(&zoneStr); err != nil {
		return nil, helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}
	loc, err := helpers.ChooseZone(c, zoneStr)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func (db *DB) UpdateUserZone(c echo.Context, userID int32, loc *time.Location) error {
	ctx := c.Request().Context()
	res, err := db.pool.Exec(ctx, updateUserZoneSQL, loc.String(), userID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if res.RowsAffected() == 0 {
		return helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	return nil
}

func (db *DB) CreateUser(c echo.Context, firstName, secondName, email string, loc *time.Location) (userID int32, err error) {
	ctx := c.Request().Context()
	row := db.pool.QueryRow(ctx, createUserSQL, firstName, secondName, email, loc.String())
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

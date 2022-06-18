package helpers

import (
	"calendar/internal/models"
	"errors"

	"github.com/labstack/echo/v4"
)

func WrapSuccess(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, &models.DataError{
		Data: data,
		Err:  nil,
	})
}

func WrapError(c echo.Context, status int, errCode string) error {
	c.JSON(status, &models.DataError{
		Data: nil,
		Err: &models.InternalError{
			Code: errCode,
		},
	})
	return errors.New(errCode)
}

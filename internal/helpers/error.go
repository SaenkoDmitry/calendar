package helpers

import (
	"calendar/internal/models"

	"github.com/labstack/echo/v4"
)

func WrapJSONError(c echo.Context, status int, msg string) error {
	return c.JSON(status, &models.CustomErr{
		Msg: msg,
	})
}

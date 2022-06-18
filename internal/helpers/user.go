package helpers

import (
	"calendar/internal/constants"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/jackc/pgconn"
)

func IsEmailDuplicated(err *pgconn.PgError) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Message, "duplicate key")
}

func GetUser(c echo.Context, paramName string) (int32, error) {
	temp, err := strconv.ParseInt(c.Param(paramName), 10, 32)
	if err != nil || temp <= 0 {
		return 0, WrapError(c, http.StatusBadRequest, constants.InvalidUserID)
	}
	return int32(temp), nil
}

package helpers

import (
	"calendar/internal/constants"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetMeeting(c echo.Context, paramName string) (int32, error) {
	temp, err := strconv.ParseInt(c.Param(paramName), 10, 32)
	if err != nil || temp <= 0 {
		return 0, WrapError(c, http.StatusBadRequest, constants.InvalidMeetingID)
	}
	return int32(temp), nil
}

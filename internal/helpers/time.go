package helpers

import (
	"calendar/internal/constants"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetDateTime(c echo.Context, input string, loc *time.Location, msgErr string) (time.Time, error) {
	parsedTime, err := time.ParseInLocation(constants.DateTimeFormat, input, loc)
	if err != nil {
		return time.Now(), WrapError(c, http.StatusBadRequest, msgErr)
	}
	return parsedTime, nil
}

func MergeDateAndTimeFromServer(date *time.Time, t *time.Time) time.Time {
	temp := date.Add(time.Hour*time.Duration(t.Hour()) + time.Minute*time.Duration(t.Minute()))
	res, _ := time.ParseInLocation(constants.DateTimeFormat, temp.Format(constants.DateTimeFormat), constants.ServerTimeZone)
	return res
}

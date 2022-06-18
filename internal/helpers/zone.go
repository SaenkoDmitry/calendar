package helpers

import (
	"calendar/internal/constants"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func ChooseZone(c echo.Context, zone string) (*time.Location, error) {
	if zone == "" {
		return time.UTC, nil
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return nil, WrapJSONError(c, http.StatusBadRequest, constants.NotValidTimeZone)
	}
	return loc, nil
}
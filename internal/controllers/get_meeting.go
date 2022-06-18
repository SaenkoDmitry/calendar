package controllers

import (
	"calendar/internal/helpers"
	"calendar/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetMeeting(c echo.Context) error {
	zone := c.QueryParam("zone")
	loc, err := helpers.ChooseZone(c, zone)
	if err != nil {
		return err
	}

	meetingID, err := helpers.GetMeeting(c, "meetingID")
	if err != nil {
		return err
	}

	meetInfo, err := repository.GetMeeting(c, h.pool, meetingID, loc)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusOK, meetInfo)
}

package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"calendar/internal/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetMeetingsByUserAndTimeInterval(c echo.Context) error {
	req := new(models.GetMeetingsByUserAndTimeRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	userID, err := helpers.GetUser(c, "userID")
	if err != nil {
		return err
	}

	loc, err := repository.SelectUserZone(c, h.pool, userID)
	if err != nil {
		return err
	}

	from, err := helpers.GetDateTime(c, req.From, loc, constants.InvalidFromDate)
	if err != nil {
		return err
	}

	to, err := helpers.GetDateTime(c, req.To, loc, constants.InvalidToDate)
	if err != nil {
		return err
	}

	if from.After(to) {
		return helpers.WrapError(c, http.StatusBadRequest, constants.FromEarlierThanToDate)
	}

	meetings, err := repository.SelectMeetingsByUserAndInterval(c, h.pool, userID, loc, from, to)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusOK, meetings)
}

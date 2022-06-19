package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetMeetingsByUserAndTimeInterval godoc
// @Summary  get meetings by user and time interval
// @Tags     meeting
// @Accept   json
// @Produce  json
// @Param    userID     path      int32  true  "User ID"
// @Param    meetingID  path      int32  true  "Meeting ID"
// @Success  200        {object}  models.DataError
// @Failure  400        {object}  models.DataError
// @Failure  500        {object}  models.DataError
// @Router   /users/{userID}/meetings [get]
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

	loc, err := h.DB.SelectUserZone(c, userID)
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

	meetings, err := h.DB.SelectMeetingsByUserAndInterval(c, userID, loc, from, to)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusOK, meetings)
}

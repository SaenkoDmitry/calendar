package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateMeeting godoc
// @Summary  create meeting info
// @Tags     meeting
// @Accept   json
// @Produce  json
// @Param    CreateMeetingReq  body      models.CreateMeetingReq  true  "request body for creating meeting"
// @Success  201               {object}  models.DataError
// @Failure  400               {object}  models.DataError
// @Failure  500               {object}  models.DataError
// @Router   /meetings [post]
func (h *handler) CreateMeeting(c echo.Context) error {
	req := new(models.CreateMeetingReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	if len(req.UserIDs) > 200 {
		return helpers.WrapError(c, http.StatusBadRequest, constants.TooManyUsersForMeeting)
	}

	loc, err := h.DB.SelectUserZone(c, req.AdminID)
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

	meetingID, err := h.DB.CreateMeetingWithLinkToUsers(c, req.AdminID, req.UserIDs, req.Name, req.Description, from, to)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusCreated, meetingID)
}

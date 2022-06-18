package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/repository"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ChangeStatusOfMeeting godoc
// @Summary  change status of meeting
// @Tags     status
// @Accept   json
// @Produce  json
// @Param    userID     path      int32   true   "User ID"
// @Param    meetingID  path      int32   true   "Meeting ID"
// @Param    status     query     string  false  "example: requested | approved | declined | finished | canceled"
// @Success  200        {object}  models.DataError
// @Failure  400        {object}  models.DataError
// @Failure  500        {object}  models.DataError
// @Router   /users/{userID}/meetings/{meetingID} [put]
func (h *handler) ChangeStatusOfMeeting(c echo.Context) error {
	userID, err := helpers.GetUser(c, "userID")
	if err != nil {
		return err
	}

	meetingID, err := helpers.GetMeeting(c, "meetingID")
	if err != nil {
		return err
	}

	newStatus := c.QueryParam("status")
	if newStatus == "" {
		return helpers.WrapError(c, http.StatusBadRequest, constants.EmptyStatus)
	}
	if !helpers.IsValidStatus(newStatus) || newStatus == constants.Requested {
		return helpers.WrapError(c, http.StatusBadRequest, constants.InvalidStatus)
	}

	currentStatus, err := repository.SelectMeetStatus(c, h.pool, userID, meetingID)
	if err != nil {
		return err
	}

	if currentStatus == constants.Canceled || currentStatus == constants.Finished {
		return helpers.WrapError(c, http.StatusBadRequest, constants.MeetingCanceledOrFinished)
	}

	err = repository.UpdateMeetStatus(c, h.pool, userID, meetingID, newStatus)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusOK, fmt.Sprintf("status changed to '%s'", newStatus))
}

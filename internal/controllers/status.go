package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) ChangeStatusOfMeeting(c echo.Context) error {
	ctx := c.Request().Context()

	temp, err := strconv.ParseInt(c.Param("userID"), 10, 32)
	if err != nil || temp <= 0 {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.InvalidUserID)
	}
	userID := int32(temp)

	temp, err = strconv.ParseInt(c.Param("meetingID"), 10, 32)
	if err != nil || temp <= 0 {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.InvalidMeetingID)
	}
	meetingID := int32(temp)

	newStatus := c.QueryParam("status")
	if newStatus == "" {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.EmptyStatus)
	}
	if !helpers.IsValidStatus(newStatus) || newStatus == constants.Requested {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.InvalidStatus)
	}

	var currentStatus string
	row := h.pool.QueryRow(ctx, "SELECT status FROM user_meetings "+
		"	WHERE user_id = $1 AND meeting_id = $2", userID, meetingID)
	if err = row.Scan(&currentStatus); err != nil {
		if err.Error() == constants.NoRowsInResultSetDBError {
			return helpers.WrapJSONError(c, http.StatusBadRequest, constants.UserNotInvitedOnTheMeeting)
		}
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if currentStatus == constants.Canceled || currentStatus == constants.Finished {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.MeetingCanceledOrFinished)
	}

	res, err := h.pool.Exec(ctx, "UPDATE user_meetings SET status = $1 WHERE user_id = $2 AND meeting_id = $3",
		newStatus, userID, meetingID)
	if err != nil {
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if res.RowsAffected() == 1 {
		return c.JSON(http.StatusOK, fmt.Sprintf("status changed to %s", newStatus))
	}

	return nil
}

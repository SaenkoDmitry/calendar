package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

func (h *handler) ChangeStatusOfMeeting(c echo.Context) error {
	ctx := c.Request().Context()

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

	var currentStatus string
	row := h.pool.QueryRow(ctx, "SELECT status FROM user_meetings "+
		"	WHERE user_id = $1 AND meeting_id = $2", userID, meetingID)
	if err = row.Scan(&currentStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return helpers.WrapError(c, http.StatusBadRequest, constants.UserNotInvitedOnTheMeeting)
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if currentStatus == constants.Canceled || currentStatus == constants.Finished {
		return helpers.WrapError(c, http.StatusBadRequest, constants.MeetingCanceledOrFinished)
	}

	_, err = h.pool.Exec(ctx, "UPDATE user_meetings SET status = $1 WHERE user_id = $2 AND meeting_id = $3",
		newStatus, userID, meetingID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return helpers.WrapSuccess(c, http.StatusOK, fmt.Sprintf("status changed to %s", newStatus))
}

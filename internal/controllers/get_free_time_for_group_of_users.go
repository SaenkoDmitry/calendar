package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	maxFetchRowsNumber = 100
	//roundingInterval   = 30
)

// GetOptimalMeetTimeForGroupOfUsers godoc
// @Summary  get first time interval for meeting for group of users when all of them are free
// @Tags     meeting
// @Accept   json
// @Produce  json
// @Param    FindOptimalMeetingTimeRequest  body      models.FindOptimalMeetingTimeRequest  true  "request body for creating meeting"
// @Success  201                            {object}  models.DataError
// @Failure  400                            {object}  models.DataError
// @Failure  500                            {object}  models.DataError
// @Router   /meetings/suggest [post]
func (h *handler) GetOptimalMeetTimeForGroupOfUsers(c echo.Context) error {
	req := new(models.FindOptimalMeetingTimeRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	existedUserIDs, err := h.DB.CheckUsersExistence(c, req.UserIDs)
	if err != nil {
		return err
	}
	if v, ok := helpers.ValidateExistenceOfUsers(req.UserIDs, existedUserIDs); !ok {
		return helpers.WrapErrorWithMsg(c, http.StatusBadRequest, constants.InvalidUserList, fmt.Sprintf("not exists next users: %v", v))
	}

	var resStart, resEnd *time.Time

	offset := 0
	//delayMinutesFromNowInterval := roundingInterval - time.Now().Minute()%roundingInterval
	prev := time.Now()
	//.Add(time.Duration(delayMinutesFromNowInterval) * time.Minute)

	for {
		// values already sorted
		meetings, err := h.DB.FindOptimalMeetingAfterCertainMoment(c, req.UserIDs, prev, maxFetchRowsNumber, offset)
		if err != nil {
			return err
		}
		if len(meetings) == 0 {
			break
		}

		for i := 0; i < len(meetings); i++ {
			fromPrevMinutesPassed := meetings[i].From.Sub(prev).Minutes()
			if fromPrevMinutesPassed > float64(req.MinDurationMinutes) {
				resStart = &prev
				temp := prev.Add(time.Minute * time.Duration(req.MinDurationMinutes))
				resEnd = &temp
				break
			}
			prev = meetings[i].To
		}

		if resStart != nil {
			break
		}
		offset += maxFetchRowsNumber
	}

	if resStart == nil { // || time.Now().Sub(*resStart).Minutes() > float64(delayMinutesFromNowInterval)
		return helpers.WrapSuccess(c, http.StatusOK, constants.NotFoundOptimalMeetingForTheInterval)
	}

	return helpers.WrapSuccess(c, http.StatusOK, &models.FindOptimalMeetingTimeResponse{
		From: resStart.Format(constants.PrettyDateTimeFormat),
		To:   resEnd.Format(constants.PrettyDateTimeFormat),
	})
}

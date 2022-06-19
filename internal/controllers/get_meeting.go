package controllers

import (
	"calendar/internal/helpers"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetMeeting godoc
// @Summary  get meeting info
// @Tags     meeting
// @Accept   json
// @Produce  json
// @Param    meetingID  path      int32  true  "Meeting ID"
// @Success  200        {object}  models.DataError
// @Failure  400        {object}  models.DataError
// @Failure  500        {object}  models.DataError
// @Router   /meetings/{meetingID} [get]
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

	meetInfo, err := h.DB.GetMeeting(c, meetingID, loc)
	if err != nil {
		return err
	}

	return helpers.WrapSuccess(c, http.StatusOK, meetInfo)
}

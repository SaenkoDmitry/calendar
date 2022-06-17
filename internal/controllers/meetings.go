package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type CreateMeetingReq struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	AdminID     int32   `json:"organizer_id"`
	UserIDs     []int32 `json:"participants"`
	From        string  `json:"from"` // 2022-01-02T15:00
	To          string  `json:"to"`   // 2022-01-02T16:00
}

func (h *handler) CreateMeeting(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(CreateMeetingReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	// extract user time zone
	zoneStr := time.UTC.String()
	rows, err := h.pool.Query(ctx, "SELECT user_zone FROM users WHERE id = $1", req.AdminID)
	if err != nil {
		return utils.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	defer rows.Close()

	if rows.Next() {
		_ = rows.Scan(&zoneStr)
	} else {
		return utils.WrapJSONError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	// validate 'from' time
	loc, err := time.LoadLocation(zoneStr)
	if err != nil {
		loc = time.UTC // set to default
	}
	from, err := time.ParseInLocation(constants.DateTimeFormat, req.From, loc)
	if err != nil {
		return utils.WrapJSONError(c, http.StatusBadRequest, constants.InvalidFromDate)
	}

	// validate 'to' time
	to, err := time.ParseInLocation(constants.DateTimeFormat, req.To, loc)
	if err != nil {
		return utils.WrapJSONError(c, http.StatusBadRequest, constants.InvalidToDate)
	}

	// validate that 'from' earlier than 'to' time
	if from.After(to) {
		return utils.WrapJSONError(c, http.StatusBadRequest, constants.FromEarlierThanToDate)
	}

	// inserting meeting and getting meeting_id
	var meetingID int32 = -1
	rows, err = h.pool.Query(ctx, "INSERT INTO meetings(meet_name, description, "+
		"	start_date, start_time, end_date, end_time) "+
		"	VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Name, req.Description,
		from.Format(constants.DateFormat),
		from.Format(constants.TimeFormat),
		to.Format(constants.DateFormat),
		to.Format(constants.TimeFormat),
	)
	if err != nil {
		return utils.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	if rows.Next() {
		_ = rows.Scan(&meetingID)
	}

	if meetingID == -1 { // for some reasons not returned meeting_id after creation
		return utils.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	// TODO append all user_meetings

	return c.JSON(http.StatusCreated, meetingID)
}

func (h *handler) GetMeeting(c echo.Context) error {

	return nil
}

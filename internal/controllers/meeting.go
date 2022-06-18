package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/labstack/echo/v4"
)

const (
	BatchInsertUserMeetingsSQL = "INSERT INTO user_meetings(user_id, meeting_id) VALUES ($1, $2)"
)

func (h *handler) CreateMeeting(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.CreateMeetingReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	// extract user time zone
	zoneStr := time.UTC.String()
	row := h.pool.QueryRow(ctx, "SELECT user_zone FROM users WHERE id = $1", req.AdminID)
	if err := row.Scan(&zoneStr); err != nil {
		return helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	loc, err := helpers.ChooseZone(c, zoneStr)
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

	// validate that 'from' earlier than 'to' time
	if from.After(to) {
		return helpers.WrapError(c, http.StatusBadRequest, constants.FromEarlierThanToDate)
	}

	tx, err := h.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
	}()
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	// inserting meeting and getting meeting_id
	var meetingID int32 = -1
	row = tx.QueryRow(ctx, "INSERT INTO meetings(meet_name, description, "+
		"	start_date, start_time, end_date, end_time) "+
		"	VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		req.Name, req.Description,
		from.In(time.UTC).Format(constants.DateFormat),
		from.In(time.UTC).Format(constants.TimeFormat),
		to.In(time.UTC).Format(constants.DateFormat),
		to.In(time.UTC).Format(constants.TimeFormat),
	)
	if err = row.Scan(&meetingID); err != nil || meetingID == -1 {
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.CannotInsertMeeting)
	}

	batch := &pgx.Batch{}

	batch.Queue(BatchInsertUserMeetingsSQL, req.AdminID, meetingID) // TODO need mark as organizer
	for _, id := range req.UserIDs {
		batch.Queue(BatchInsertUserMeetingsSQL, id, meetingID)
	}
	batchResults := tx.SendBatch(ctx, batch)
	err = batchResults.Close()
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if pgErr.Message == constants.UserIDConstraintDBErr {
			return helpers.WrapError(c, http.StatusInternalServerError, pgErr.Detail)
		}
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	if err = tx.Commit(ctx); err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return helpers.WrapSuccess(c, http.StatusCreated, meetingID)
}

func (h *handler) GetMeeting(c echo.Context) error {
	ctx := c.Request().Context()
	zone := c.QueryParam("zone")
	loc, err := helpers.ChooseZone(c, zone)
	if err != nil {
		return err
	}

	meetingID, err := helpers.GetMeeting(c, "id")
	if err != nil {
		return err
	}

	var meetName, description string
	var startDate, startTime, endDate, endTime *time.Time

	row := h.pool.QueryRow(ctx, "SELECT meet_name, description, start_date, start_time, end_date, end_time "+
		"	FROM meetings WHERE id = $1", meetingID)
	if err := row.Scan(&meetName, &description, &startDate, &startTime, &endDate, &endTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return helpers.WrapError(c, http.StatusBadRequest, constants.MeetingIDNotExists)
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	fromTime := helpers.MergeDateAndTime(startDate, startTime)
	toTime := helpers.MergeDateAndTime(endDate, endTime)

	meetInfo := &models.MeetingInfoResponse{
		Name:        meetName,
		Description: description,
		From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
		To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
	}

	return helpers.WrapSuccess(c, http.StatusOK, meetInfo)
}

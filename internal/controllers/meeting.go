package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"net/http"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/labstack/echo/v4"
)

const (
	UserIDConstraint           = `insert or update on table "user_meetings" violates foreign key constraint "fk_user_id"`
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
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	// validate 'from' time
	loc, err := time.LoadLocation(zoneStr)
	if err != nil {
		loc = time.UTC // set to default
	}
	from, err := time.ParseInLocation(constants.DateTimeFormat, req.From, loc)
	if err != nil {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.InvalidFromDate)
	}

	// validate 'to' time
	to, err := time.ParseInLocation(constants.DateTimeFormat, req.To, loc)
	if err != nil {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.InvalidToDate)
	}

	// validate that 'from' earlier than 'to' time
	if from.After(to) {
		return helpers.WrapJSONError(c, http.StatusBadRequest, constants.FromEarlierThanToDate)
	}

	tx, err := h.pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
	}()
	if err != nil {
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
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
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.CannotInsertMeeting)
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
		if pgErr.Message == UserIDConstraint {
			return helpers.WrapJSONError(c, http.StatusInternalServerError, pgErr.Detail)
		}
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	if err = tx.Commit(ctx); err != nil {
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return c.JSON(http.StatusCreated, meetingID)
}

func (h *handler) GetMeeting(c echo.Context) error {
	ctx := c.Request().Context()
	meetingID := c.Param("id")
	zone := c.QueryParam("zone")

	loc, err := helpers.ChooseZone(c, zone)
	if err != nil {
		return err
	}

	var meetName, description string
	var startDate, startTime, endDate, endTime *time.Time

	row := h.pool.QueryRow(ctx, "SELECT meet_name, description, start_date, start_time, end_date, end_time "+
		"	FROM meetings WHERE id = $1", meetingID)
	if err := row.Scan(&meetName, &description, &startDate, &startTime, &endDate, &endTime); err != nil {
		if err.Error() == constants.NoRowsInResultSetDBError {
			return helpers.WrapJSONError(c, http.StatusBadRequest, constants.MeetingIDNotExists)
		}
		return helpers.WrapJSONError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	fromTime := startDate.Add(time.Hour*time.Duration(startTime.Hour()) + time.Minute*time.Duration(startTime.Minute()))
	toTime := endDate.Add(time.Hour*time.Duration(endTime.Hour()) + time.Minute*time.Duration(endTime.Minute()))

	meetInfo := &models.MeetingInfoResponse{
		Name:        meetName,
		Description: description,
		From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
		To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
	}

	return c.JSON(http.StatusOK, meetInfo)
}

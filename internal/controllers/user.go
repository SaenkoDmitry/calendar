package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgconn"

	"github.com/labstack/echo/v4"
)

func (h *handler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.CreateUserReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	loc, err := helpers.ChooseZone(c, req.Zone)
	if err != nil {
		return err
	}

	_, err = h.pool.Exec(ctx, "INSERT INTO users(first_name, second_name, email) "+
		"	VALUES($1, $2, $3)", req.FirstName, req.SecondName, req.Email, loc.String())
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if helpers.IsEmailDuplicated(pgErr) {
			return helpers.WrapError(c, http.StatusConflict, constants.EmailAlreadyRegistered)
		}
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return helpers.WrapSuccess(c, http.StatusCreated, req.Email)
}

func (h *handler) ChangeUserZone(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.ChangeZoneReq)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	userID, err := helpers.GetUser(c, "id")
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation(req.Zone)
	if err != nil {
		return helpers.WrapError(c, http.StatusBadRequest, constants.NotValidTimeZone)
	}

	res, err := h.pool.Exec(ctx, "UPDATE users SET user_zone = $1 WHERE id = $2", loc.String(), userID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	if res.RowsAffected() == 0 {
		return helpers.WrapError(c, http.StatusBadRequest, constants.UserIDNotExists)
	}

	_ = res

	return helpers.WrapSuccess(c, http.StatusOK, loc)
}

func (h *handler) GetMeetingsByUserAndTimeInterval(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.GetMeetingsByUserAndTimeRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	userID, err := helpers.GetUser(c, "id")
	if err != nil {
		return err
	}

	zoneStr := time.UTC.String()
	row := h.pool.QueryRow(ctx, "SELECT user_zone FROM users WHERE id = $1", userID)
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

	rows, err := h.pool.Query(ctx, `SELECT 
			m.meet_name, m.description, m.start_date, m.start_time, m.end_date, m.end_time
		FROM user_meetings um 
				JOIN meetings m ON m.id = um.meeting_id 
		WHERE um.user_id = $1 AND um.status != 'canceled' AND um.status != 'finished' 
			AND
			(m.start_date > $2 OR m.start_date = $2 AND m.start_time >= $3) 
			AND 
			(m.start_date < $4 OR m.start_date = $4 AND m.start_time < $5)`,
		userID,
		from.In(time.UTC).Format(constants.DateFormat),
		from.In(time.UTC).Format(constants.TimeFormat),
		to.In(time.UTC).Format(constants.DateFormat),
		to.In(time.UTC).Format(constants.TimeFormat),
	)
	defer rows.Close()
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	meetings := make([]*models.MeetingInfoResponse, 0)
	for rows.Next() {
		var meetName, description string
		var startDate, startTime, endDate, endTime *time.Time
		if err = rows.Scan(&meetName, &description, &startDate, &startTime, &endDate, &endTime); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return helpers.WrapError(c, http.StatusBadRequest, constants.MeetingIDNotExists)
			}
			return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
		}
		fromTime := helpers.MergeDateAndTime(startDate, startTime)
		toTime := helpers.MergeDateAndTime(endDate, endTime)
		meetings = append(meetings, &models.MeetingInfoResponse{
			Name:        meetName,
			Description: description,
			From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
			To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
		})
	}

	return helpers.WrapSuccess(c, http.StatusOK, meetings)
}

func (h *handler) GetFreeTimeForGroupOfUsers(c echo.Context) error {

	return nil
}

package repository

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"errors"
	"net/http"
	"time"

	"github.com/lib/pq"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	selectMeetingsByUserAndIntervalSQL = `SELECT 
			m.id, m.meet_name, m.description, m.start_date, m.start_time, m.end_date, m.end_time
		FROM user_meetings um 
				JOIN meetings m ON m.id = um.meeting_id 
		WHERE um.user_id = $1 AND um.status != 'canceled' AND um.status != 'finished' 
			AND
			(m.start_date > $2 OR m.start_date = $2 AND m.start_time >= $3) 
			AND 
			(m.start_date < $4 OR m.start_date = $4 AND m.start_time < $5)`

	selectMeetingStatusSQL = `SELECT status FROM user_meetings 
		WHERE user_id = $1 AND meeting_id = $2`

	updateMeetStatusSQL = `UPDATE user_meetings SET status = $1 WHERE user_id = $2 AND meeting_id = $3`

	insertMeetingSQL = `INSERT INTO meetings(meet_name, description, 
		start_date, start_time, end_date, end_time) 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	batchInsertUserMeetingsSQL = `INSERT INTO user_meetings(user_id, meeting_id) VALUES ($1, $2)`
	getMeetingSQL              = `SELECT id, meet_name, description, start_date, start_time, end_date, end_time FROM meetings WHERE id = $1`

	selectFirstNMeetingsByUserGroupSQL = `SELECT meeting_id, start_date, start_time, end_date, end_time
		FROM user_meetings um
			JOIN meetings m on um.meeting_id = m.id
		WHERE um.status != 'canceled' AND um.status != 'finished' AND um.status != 'declined'
  			AND m.start_date >= $1 AND m.start_time >= $2
  			AND um.user_id = ANY($3)
		GROUP BY meeting_id, start_date, start_time, end_date, end_time
		ORDER BY start_date, start_time, end_date, end_time ASC
		LIMIT $4 OFFSET $5`
)

func SelectMeetingsByUserAndInterval(c echo.Context, pool *pgxpool.Pool,
	userID int32, loc *time.Location, from time.Time, to time.Time) ([]*models.MeetingInfoResponse, error) {
	ctx := c.Request().Context()
	rows, err := pool.Query(ctx, selectMeetingsByUserAndIntervalSQL, userID,
		from.In(constants.ServerTimeZone).Format(constants.DateFormat),
		from.In(constants.ServerTimeZone).Format(constants.TimeFormat),
		to.In(constants.ServerTimeZone).Format(constants.DateFormat),
		to.In(constants.ServerTimeZone).Format(constants.TimeFormat),
	)
	defer rows.Close()
	if err != nil {
		return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	meetings := make([]*models.MeetingInfoResponse, 0)
	for rows.Next() {
		var ID int32
		var meetName, description string
		var startDate, startTime, endDate, endTime *time.Time
		if err = rows.Scan(&ID, &meetName, &description, &startDate, &startTime, &endDate, &endTime); err != nil {
			return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
		}
		fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
		toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)
		meetings = append(meetings, &models.MeetingInfoResponse{
			ID:          ID,
			Name:        meetName,
			Description: description,
			From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
			To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
		})
	}
	return meetings, nil
}

func SelectMeetStatus(c echo.Context, pool *pgxpool.Pool, userID, meetingID int32) (string, error) {
	ctx := c.Request().Context()
	var status string
	row := pool.QueryRow(ctx, selectMeetingStatusSQL, userID, meetingID)
	if err := row.Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", helpers.WrapError(c, http.StatusBadRequest, constants.UserNotInvitedOnTheMeeting)
		}
		return "", helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	return status, nil
}

func UpdateMeetStatus(c echo.Context, pool *pgxpool.Pool, userID, meetingID int32, status string) error {
	ctx := c.Request().Context()
	_, err := pool.Exec(ctx, updateMeetStatusSQL,
		status, userID, meetingID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	return nil
}

func CreateMeetingWithLinkToUsers(c echo.Context, pool *pgxpool.Pool,
	adminID int32, userIDs []int32,
	name, description string, from, to time.Time,
) (int32, error) {
	ctx := c.Request().Context()
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	defer func() {
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
	}()
	if err != nil {
		return 0, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	var meetingID int32 = -1
	row := tx.QueryRow(ctx, insertMeetingSQL,
		name, description,
		from.In(constants.ServerTimeZone).Format(constants.DateFormat),
		from.In(constants.ServerTimeZone).Format(constants.TimeFormat),
		to.In(constants.ServerTimeZone).Format(constants.DateFormat),
		to.In(constants.ServerTimeZone).Format(constants.TimeFormat),
	)
	if err = row.Scan(&meetingID); err != nil || meetingID == -1 {
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
		return 0, helpers.WrapError(c, http.StatusInternalServerError, constants.CannotCreateMeeting)
	}

	batch := &pgx.Batch{}

	batch.Queue(batchInsertUserMeetingsSQL, adminID, meetingID) // TODO need mark as organizer
	for _, id := range userIDs {
		batch.Queue(batchInsertUserMeetingsSQL, id, meetingID)
	}
	batchResults := tx.SendBatch(ctx, batch)
	err = batchResults.Close()
	if err != nil {
		pgErr, _ := err.(*pgconn.PgError)
		if pgErr.Message == constants.UserIDConstraintDBErr {
			return 0, helpers.WrapError(c, http.StatusInternalServerError, pgErr.Detail)
		}
		if e := tx.Rollback(ctx); e != nil {
			c.Logger().Error(e)
		}
		return 0, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	return meetingID, nil
}

func GetMeeting(c echo.Context, pool *pgxpool.Pool, meetingID int32, loc *time.Location) (*models.MeetingInfoResponse, error) {
	ctx := c.Request().Context()
	var ID int32
	var meetName, description string
	var startDate, startTime, endDate, endTime *time.Time

	row := pool.QueryRow(ctx, getMeetingSQL, meetingID)
	if err := row.Scan(&ID, &meetName, &description, &startDate, &startTime, &endDate, &endTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helpers.WrapError(c, http.StatusBadRequest, constants.MeetingIDNotExists)
		}
		return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
	toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)

	return &models.MeetingInfoResponse{
		ID:          ID,
		Name:        meetName,
		Description: description,
		From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
		To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
	}, nil
}

func FindOptimalMeetingAfterCertainMoment(c echo.Context, pool *pgxpool.Pool, userIDs []int32,
	startingPoint time.Time, count, offset int) ([]*models.MeetingDataForOptimalCalcTime, error) {
	ctx := c.Request().Context()
	rows, err := pool.Query(ctx, selectFirstNMeetingsByUserGroupSQL,
		startingPoint.Format(constants.DateFormat),
		startingPoint.Format(constants.TimeFormat),
		pq.Array(userIDs),
		count, offset, // to do not fetch too long values immediately and not get very slow queries
	)
	defer rows.Close()
	if err != nil {
		return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	meetings := make([]*models.MeetingDataForOptimalCalcTime, 0)
	for rows.Next() {
		var ID int32
		var startDate, startTime, endDate, endTime *time.Time
		err = rows.Scan(&ID, &startDate, &startTime, &endDate, &endTime)
		if err != nil {
			pgErr, _ := err.(*pgconn.PgError)
			_ = pgErr
		}
		fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
		toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)
		meetings = append(meetings, &models.MeetingDataForOptimalCalcTime{
			ID:   ID,
			From: fromTime,
			To:   toTime,
		})
	}
	return meetings, err
}

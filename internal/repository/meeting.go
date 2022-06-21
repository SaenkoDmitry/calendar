package repository

import (
	"calendar/internal/constants"
	"calendar/internal/helpers"
	"calendar/internal/models"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgtype"

	"github.com/lib/pq"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

const (
	selectMeetingsByUserAndIntervalSQL = `SELECT 
			m.id, m.meet_name, m.description, m.repeat, m.start_date, m.start_time, m.end_date, m.end_time
		FROM user_meetings um 
				JOIN meetings m ON m.id = um.meeting_id 
		WHERE um.user_id = $1 
			AND
			(m.start_date > $2 OR m.start_date = $2 AND m.start_time >= $3) 
			AND 
			(m.start_date < $4 OR m.start_date = $4 AND m.start_time < $5)`

	selectVirtualMeetingsByUserBeforeToSQL = `SELECT 
			m.id, m.repeat, m.start_date, m.start_time, m.end_date, m.end_time
		FROM user_meetings um 
				JOIN meetings m ON m.id = um.meeting_id 
		WHERE um.user_id = $1 AND repeat IS NOT NULL AND
			(m.start_date < $2 OR m.start_date = $2 AND m.start_time <= $3)`

	selectMeetingsByIDsSQL = `SELECT 
			m.id, m.meet_name, m.description, m.repeat, m.start_date, m.start_time, m.end_date, m.end_time
		FROM meetings m
		WHERE m.id = ANY($1)`

	updateMeetStatusSQL = `UPDATE user_meetings SET status = $1 WHERE user_id = $2 AND meeting_id = $3`

	insertMeetingSQL = `INSERT INTO meetings(meet_name, description, repeat, 
		start_date, start_time, end_date, end_time) 
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	batchInsertUserMeetingsSQL = `INSERT INTO user_meetings(user_id, meeting_id, status) VALUES ($1, $2, $3)`
	getMeetingSQL              = `SELECT m.id as id, array_agg(u.id) as participants, meet_name, description, start_date, start_time, end_date, end_time
    	FROM meetings m
    		LEFT JOIN user_meetings um on m.id = um.meeting_id
    		LEFT JOIN users u on um.user_id = u.id
    	WHERE m.id = $1
    	GROUP BY m.id
`

	selectFirstAllowedTimeIntervalByUserGroupSQL = `SELECT meeting_id, start_date, start_time, end_date, end_time
		FROM user_meetings um
			JOIN meetings m on um.meeting_id = m.id
		WHERE um.status != 'declined' AND um.user_id = ANY($1)
			AND (m.start_date > $2 OR m.start_date = $2 AND m.start_time >= $3)
		GROUP BY meeting_id, start_date, start_time, end_date, end_time
		ORDER BY start_date, start_time, end_date, end_time ASC
		LIMIT $4 OFFSET $5`
)

func (db *DB) SelectMeetingsByUserAndInterval(c echo.Context,
	userID int32, loc *time.Location, from time.Time, to time.Time) ([]*models.MeetingInfoResponse, error) {
	ctx := c.Request().Context()
	rows, err := db.pool.Query(ctx, selectMeetingsByUserAndIntervalSQL, userID,
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
		var meetName string
		var description, repeat sql.NullString
		var startDate, startTime, endDate, endTime *time.Time
		if err = rows.Scan(&ID, &meetName, &description, &repeat, &startDate, &startTime, &endDate, &endTime); err != nil {
			return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
		}
		fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
		toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)
		meetings = append(meetings, &models.MeetingInfoResponse{
			ID:          ID,
			Name:        meetName,
			Description: description.String,
			From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
			To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
			Repeat:      repeat.String,
		})
	}
	return meetings, nil
}

func (db *DB) SelectVirtualMeetingsByUserStartsBefore(c echo.Context, userID int32, to time.Time) ([]*models.VirtualMeetingInfo, error) {
	ctx := c.Request().Context()
	rows, err := db.pool.Query(ctx, selectVirtualMeetingsByUserBeforeToSQL, userID,
		to.In(constants.ServerTimeZone).Format(constants.DateFormat),
		to.In(constants.ServerTimeZone).Format(constants.TimeFormat),
	)
	defer rows.Close()
	if err != nil {
		return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	meetings := make([]*models.VirtualMeetingInfo, 0)
	for rows.Next() {
		var ID int32
		var repeat sql.NullString
		var startDate, startTime, endDate, endTime *time.Time
		if err = rows.Scan(&ID, &repeat, &startDate, &startTime, &endDate, &endTime); err != nil {
			return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
		}
		fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
		toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)
		meetings = append(meetings, &models.VirtualMeetingInfo{
			ID:     ID,
			From:   fromTime,
			To:     toTime,
			Repeat: repeat.String,
		})
	}
	return meetings, nil
}

func (db *DB) ResolveMeetingsByIDs(c echo.Context, IDs []int32, loc *time.Location) ([]*models.MeetingInfoResponse, error) {
	ctx := c.Request().Context()
	rows, err := db.pool.Query(ctx, selectMeetingsByIDsSQL, IDs)
	defer rows.Close()
	if err != nil {
		return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	meetings := make([]*models.MeetingInfoResponse, 0)
	for rows.Next() {
		var ID int32
		var meetName string
		var description, repeat sql.NullString
		var startDate, startTime, endDate, endTime *time.Time
		if err = rows.Scan(&ID, &meetName, &description, &repeat, &startDate, &startTime, &endDate, &endTime); err != nil {
			return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
		}
		fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
		toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)
		meetings = append(meetings, &models.MeetingInfoResponse{
			ID:          ID,
			Name:        meetName,
			Description: description.String,
			From:        fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
			To:          toTime.In(loc).Format(constants.PrettyDateTimeFormat),
			Repeat:      repeat.String,
		})
	}
	return meetings, nil
}

func (db *DB) UpdateMeetStatus(c echo.Context, userID, meetingID int32, status string) error {
	ctx := c.Request().Context()
	res, err := db.pool.Exec(ctx, updateMeetStatusSQL,
		status, userID, meetingID)
	if err != nil {
		return helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}
	if res.RowsAffected() == 0 {
		return helpers.WrapError(c, http.StatusNotFound, constants.NothingUpdated)
	}
	return nil
}

func (db *DB) CreateMeetingWithLinkToUsers(c echo.Context, adminID int32, userIDs []int32, name, description, repeat string, from, to time.Time) (int32, error) {
	ctx := c.Request().Context()
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
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
		name, description, sql.NullString{String: repeat, Valid: repeat != ""},
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

	batch.Queue(batchInsertUserMeetingsSQL, adminID, meetingID, constants.Approved) // TODO need mark as organizer?
	for _, id := range userIDs {
		if id == adminID {
			continue
		}
		batch.Queue(batchInsertUserMeetingsSQL, id, meetingID, constants.Requested)
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

func (db *DB) GetMeeting(c echo.Context, meetingID int32, loc *time.Location) (*models.MeetingInfoResponse, error) {
	ctx := c.Request().Context()
	var ID int32
	var participants pgtype.Int4Array
	var meetName, description string
	var startDate, startTime, endDate, endTime *time.Time

	row := db.pool.QueryRow(ctx, getMeetingSQL, meetingID)
	if err := row.Scan(&ID, &participants, &meetName, &description, &startDate, &startTime, &endDate, &endTime); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, helpers.WrapError(c, http.StatusBadRequest, constants.MeetingIDNotExists)
		}
		return nil, helpers.WrapError(c, http.StatusInternalServerError, constants.UndefinedDB)
	}

	fromTime := helpers.MergeDateAndTimeFromServer(startDate, startTime)
	toTime := helpers.MergeDateAndTimeFromServer(endDate, endTime)

	return &models.MeetingInfoResponse{
		ID:           ID,
		Name:         meetName,
		Description:  description,
		Participants: helpers.ConvertInt4ArrayToInt32(participants),
		From:         fromTime.In(loc).Format(constants.PrettyDateTimeFormat),
		To:           toTime.In(loc).Format(constants.PrettyDateTimeFormat),
	}, nil
}

func (db *DB) FindOptimalMeetingAfterCertainMoment(c echo.Context, userIDs []int32,
	startingPoint time.Time, count, offset int) ([]*models.MeetingDataForOptimalCalcTime, error) {
	ctx := c.Request().Context()
	rows, err := db.pool.Query(ctx, selectFirstAllowedTimeIntervalByUserGroupSQL,
		pq.Array(userIDs),
		startingPoint.Format(constants.DateFormat),
		startingPoint.Format(constants.TimeFormat),
		count,
		offset, // to do not fetch too long values immediately and not get very slow queries
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
	return meetings, nil
}

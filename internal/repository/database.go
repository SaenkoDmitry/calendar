package repository

import (
	"calendar/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type DBService interface {
	Check(ctx context.Context) string

	SelectMeetingsByUserAndInterval(c echo.Context, userID int32, loc *time.Location,
		from time.Time, to time.Time) ([]*models.MeetingInfoResponse, error)

	UpdateMeetStatus(c echo.Context, userID, meetingID int32, status string) error

	CreateMeetingWithLinkToUsers(c echo.Context,
		adminID int32, userIDs []int32,
		name, description string, from, to time.Time,
	) (int32, error)

	GetMeeting(c echo.Context, meetingID int32, loc *time.Location) (*models.MeetingInfoResponse, error)

	FindOptimalMeetingAfterCertainMoment(c echo.Context, userIDs []int32,
		startingPoint time.Time, count, offset int) ([]*models.MeetingDataForOptimalCalcTime, error)

	SelectUserZone(c echo.Context, userID int32) (*time.Location, error)

	UpdateUserZone(c echo.Context, userID int32, loc *time.Location) error

	CreateUser(c echo.Context, firstName, secondName, email string, loc *time.Location) (userID int32, err error)
}

type DB struct {
	pool *pgxpool.Pool
}

func NewDBService(pool *pgxpool.Pool) DBService {
	return &DB{
		pool: pool,
	}
}

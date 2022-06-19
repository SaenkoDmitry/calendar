package controllers

import (
	"calendar/internal/models"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func (m MockDB) SelectMeetingsByUserAndInterval(c echo.Context, userID int32, loc *time.Location, from time.Time, to time.Time) ([]*models.MeetingInfoResponse, error) {
	return []*models.MeetingInfoResponse{}, nil
}

func Test_handler_GetMeetingsByUserAndTimeInterval(t *testing.T) {
	// TODO Test_handler_GetMeetingsByUserAndTimeInterval
}

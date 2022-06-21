package controllers

import (
	"calendar/internal/models"
	"time"

	"github.com/labstack/echo/v4"
)

type MockDB struct {
}

func (m MockDB) SelectVirtualMeetingsByUserStartsBefore(c echo.Context, userID int32, to time.Time) ([]*models.VirtualMeetingInfo, error) {
	return []*models.VirtualMeetingInfo{}, nil
}

func (m MockDB) ResolveMeetingsByIDs(c echo.Context, IDs []int32, loc *time.Location) ([]*models.MeetingInfoResponse, error) {
	return []*models.MeetingInfoResponse{}, nil
}

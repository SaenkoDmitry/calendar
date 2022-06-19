package controllers

import (
	"calendar/internal/models"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func (m MockDB) GetMeeting(c echo.Context, meetingID int32, loc *time.Location) (*models.MeetingInfoResponse, error) {
	return &models.MeetingInfoResponse{}, nil
}

func Test_handler_GetMeeting(t *testing.T) {
	// TODO Test_handler_GetMeeting
}

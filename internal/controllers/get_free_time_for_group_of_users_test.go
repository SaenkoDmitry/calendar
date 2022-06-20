package controllers

import (
	"calendar/internal/models"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func (m MockDB) CheckUsersExistence(c echo.Context, userIDs []int32) ([]int32, error) {
	return []int32{}, nil
}

func (m MockDB) FindOptimalMeetingAfterCertainMoment(c echo.Context, userIDs []int32, startingPoint time.Time, count, offset int) ([]*models.MeetingDataForOptimalCalcTime, error) {
	return []*models.MeetingDataForOptimalCalcTime{}, nil
}

func Test_handler_GetOptimalMeetTimeForGroupOfUsers(t *testing.T) {
	// TODO Test_handler_GetOptimalMeetTimeForGroupOfUsers
}

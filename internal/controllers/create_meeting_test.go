package controllers

import (
	"calendar/internal/middlewares"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	createMeetingReq = "{\n\t\"name\": \"Встреча\",\n\t\"organizer_id\": 8,\n\t\"participants\": [85, 87, 8],\n\t\"from\": \"2022-01-02T18:00\",\n\t\"to\": \"2022-01-02T19:00\"\n}"
)

func (m MockDB) CreateMeetingWithLinkToUsers(c echo.Context, adminID int32, userIDs []int32, name, description, repeat string, from, to time.Time) (int32, error) {
	return 5, nil
}

func (m MockDB) SelectUserZone(c echo.Context, userID int32) (*time.Location, error) {
	loc, _ := time.LoadLocation("Europe/Moscow")
	return loc, nil
}

func Test_handler_CreateMeeting_Success(t *testing.T) {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createMeetingReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/meetings")
	mockDB := &MockDB{}
	h := NewHandler(mockDB)
	if assert.NoError(t, h.CreateMeeting(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

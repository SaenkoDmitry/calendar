package controllers

import (
	"calendar/internal/constants"
	"calendar/internal/middlewares"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func (m MockDB) UpdateMeetStatus(c echo.Context, userID, meetingID int32, status string) error {
	return nil
}

func Test_handler_ChangeStatusOfMeeting_Success(t *testing.T) {
	statusChangedJSON := "{\"data\":\"status changed to 'declined'\",\"err\":null}\n"
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
	q := make(url.Values)
	q.Set("status", constants.Declined)
	req := httptest.NewRequest(http.MethodPut, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:userID/meetings/:meetingID")
	c.SetParamNames("userID", "meetingID", "status")
	c.SetParamValues("1", "37", "declined")
	mockDB := &MockDB{}
	h := NewHandler(mockDB)
	if assert.NoError(t, h.ChangeStatusOfMeeting(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, statusChangedJSON, rec.Body.String())
	}
}

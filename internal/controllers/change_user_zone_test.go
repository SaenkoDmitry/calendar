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
	validZoneJSONReq            = "{\n\t\"zone\": \"Europe/Copenhagen\"\n}"
	invalidZoneJSONReq          = "{\n\t\"zone\": \"test\"\n}"
	successfullyChangedZoneJSON = "{\"data\":\"Europe/Copenhagen\",\"err\":null}\n"
	invalidZoneJSON             = "{\"data\":null,\"err\":{\"code\":\"invalid_time_zone\",\"msg\":\"\"}}\n"
	invalidUserIDJSON           = "{\"data\":null,\"err\":{\"code\":\"invalid_user_id\",\"msg\":\"\"}}\n"
)

func (m MockDB) UpdateUserZone(c echo.Context, userID int32, loc *time.Location) error {
	return nil
}

func Test_handler_ChangeUserZone_Success(t *testing.T) {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(validZoneJSONReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:userID")
	c.SetParamNames("userID")
	c.SetParamValues("1")
	mockDB := &MockDB{}
	h := NewHandler(mockDB)
	if assert.NoError(t, h.ChangeUserZone(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, successfullyChangedZoneJSON, rec.Body.String())
	}
}

func Test_handler_ChangeUserZone_InvalidZone(t *testing.T) {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(invalidZoneJSONReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:userID")
	c.SetParamNames("userID")
	c.SetParamValues("1")
	mockDB := &MockDB{}
	h := NewHandler(mockDB)
	h.ChangeUserZone(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, invalidZoneJSON, rec.Body.String())
}

func Test_handler_ChangeUserZone_InvalidUserID(t *testing.T) {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(validZoneJSONReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:userID")
	c.SetParamNames("userID")
	c.SetParamValues("test")
	mockDB := &MockDB{}
	h := NewHandler(mockDB)
	h.ChangeUserZone(c)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, invalidUserIDJSON, rec.Body.String())
}

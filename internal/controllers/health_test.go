package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func (m MockDB) Check(ctx context.Context) string {
	return "UP"
}

func Test_handler_HealthCheck(t *testing.T) {
	postgresStatusJSON := "{\"name\":\"postgres_db\",\"status\":\"UP\"}"
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/health")
	mockDB := &MockDB{}
	h := NewHandler(mockDB)
	if assert.NoError(t, h.HealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		res := rec.Body.String()
		assert.Contains(t, res, postgresStatusJSON)
	}
}

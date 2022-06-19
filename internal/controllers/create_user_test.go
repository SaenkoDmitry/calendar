package controllers

import (
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func (m MockDB) CreateUser(c echo.Context, firstName, secondName, email string, loc *time.Location) (userID int32, err error) {
	//TODO implement me
	panic("implement me")
}

func Test_handler_CreateUser(t *testing.T) {
	// TODO Test_handler_CreateUser
}

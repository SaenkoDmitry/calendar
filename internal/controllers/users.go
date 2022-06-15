package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (h *handler) CreateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if err := c.Validate(u); err != nil {
		return err
	}

	// TODO save to db

	return c.JSON(http.StatusCreated, u)
}

func (h *handler) GetMeetingsByUser(c echo.Context) error {

	return nil
}

func (h *handler) GetMeetingsForGroupOfUsers(c echo.Context) error {

	return nil
}

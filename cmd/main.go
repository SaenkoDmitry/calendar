package main

import (
	"calendar/internal/controllers"
	"calendar/internal/middlewares"
	"net/http"

	"github.com/go-playground/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := controllers.NewHandler()

	e.POST("/users", h.CreateUser)       // создать пользователя
	e.POST("/meetings", h.CreateMeeting) // создать встречу в календаре пользователя со списком приглашенных пользователей
	e.GET("/meetings/:id", h.GetMeeting) // получить детали встречи

	// statuses of meeting: created, approved, rejected, completed
	e.PUT("/meetings/:id/status", h.ChangeStatusOfMeeting) // принять или отклонить приглашение другого пользователя
	e.GET("/users/:id/meetings", h.GetMeetingsByUser)      // найти все встречи пользователя для заданного промежутка времени
	e.GET("/meetings", h.GetMeetingsForGroupOfUsers)       // найти ближайшей интервал времени, в котором все эти пользователи свободны
	e.GET("/health", h.HealthCheck)

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

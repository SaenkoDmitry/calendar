package main

import (
	"calendar/internal/controllers"
	"calendar/internal/db"
	"calendar/internal/middlewares"
	"context"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := db.MakeMigrations()
	if err != nil {
		log.Fatalf("Could not make migrations: %s", err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	dbPool, err := db.InitPool(ctx)
	if err != nil {
		log.Fatalf("Could not set up database: %s", err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	h := controllers.NewHandler(dbPool)

	e.POST("/users", h.CreateUser)        // создать пользователя
	e.PUT("/users/:id", h.ChangeUserZone) // сменить часовой пояс пользователя
	e.POST("/meetings", h.CreateMeeting)  // создать встречу в календаре пользователя со списком приглашенных пользователей
	e.GET("/meetings/:id", h.GetMeeting)  // получить детали встречи

	// statuses of meeting: created, approved, rejected, completed
	e.PUT("/meetings/:id/status", h.ChangeStatusOfMeeting) // принять или отклонить приглашение другого пользователя
	e.GET("/users/:id/meetings", h.GetMeetingsByUser)      // найти все встречи пользователя для заданного промежутка времени
	e.GET("/meetings", h.GetMeetingsForGroupOfUsers)       // найти ближайшей интервал времени, в котором все эти пользователи свободны
	e.GET("/health", h.HealthCheck)

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

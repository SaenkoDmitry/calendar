package main

import (
	"calendar/internal/controllers"
	"calendar/internal/db"
	"calendar/internal/middlewares"
	"calendar/internal/repository"
	"net/http"
	"os"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"

	_ "calendar/docs" // docs is generated by Swag CLI

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// @title           Calendar API
// @version         1.0
// @description     calendar http server documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   Dmitry Saenko
// @contact.url    https://github.com/SaenkoDmitry
// @contact.email  dmitryssaenko@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := db.MakeMigrations()
	if err != nil {
		log.Fatalf("Could not make migrations: %s", err.Error())
		os.Exit(1)
	}

	dbPool, err := db.InitPool()
	if err != nil {
		log.Fatalf("Could not set up database: %s", err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	database := repository.NewDBService(dbPool)

	h := controllers.NewHandler(database)

	e.POST("/users", h.CreateUser)              // создать пользователя
	e.PUT("/users/:userID", h.ChangeUserZone)   // сменить часовой пояс пользователя
	e.POST("/meetings", h.CreateMeeting)        // создать встречу в календаре пользователя со списком приглашенных пользователей
	e.GET("/meetings/:meetingID", h.GetMeeting) // получить детали встречи

	e.PUT("/users/:userID/meetings/:meetingID", h.ChangeStatusOfMeeting) // принять или отклонить приглашение другого пользователя
	e.GET("/users/:userID/meetings", h.GetMeetingsByUserAndTimeInterval) // найти все встречи пользователя для заданного промежутка времени
	e.POST("/meetings/suggest", h.GetOptimalMeetTimeForGroupOfUsers)     // найти ближайшей интервал времени, в котором все эти пользователи свободны
	e.GET("/health", h.HealthCheck)

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

package controllers

import (
	"calendar/internal/repository"
)

type handler struct {
	DB repository.DBService
}

func NewHandler(database repository.DBService) *handler {
	return &handler{
		DB: database,
	}
}

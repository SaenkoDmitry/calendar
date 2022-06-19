package controllers

import (
	"calendar/internal/repository"
)

type handler struct {
	DB *repository.DB
}

func NewHandler(database *repository.DB) *handler {
	return &handler{
		DB: database,
	}
}

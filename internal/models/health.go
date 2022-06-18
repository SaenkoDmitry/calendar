package models

import "time"

type HealthStatus struct {
	Status    string           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Services  []*HealthService `json:"services"`
}

type HealthService struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

package models

import "time"

type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	DB        string    `json:"database"`
}

package models

import "time"

type CreateMeetingReq struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	AdminID     int32   `json:"organizer_id" validate:"required"`
	UserIDs     []int32 `json:"participants" validate:"required"`
	From        string  `json:"from" validate:"required"` // 2022-01-02T15:00
	To          string  `json:"to" validate:"required"`   // 2022-01-02T16:00
}

type MeetingInfoResponse struct {
	ID           int32   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Participants []int32 `json:"participants"`
	From         string  `json:"from"`
	To           string  `json:"to"`
}

type FindOptimalMeetingTimeRequest struct {
	UserIDs            []int32 `json:"participants" validate:"required"`
	MinDurationMinutes int32   `json:"min_duration_in_minutes" validate:"required"`
}

type MeetingDataForOptimalCalcTime struct {
	ID   int32
	From time.Time
	To   time.Time
}

type FindOptimalMeetingTimeResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

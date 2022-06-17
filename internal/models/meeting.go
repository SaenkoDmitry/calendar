package models

type CreateMeetingReq struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	AdminID     int32   `json:"organizer_id"`
	UserIDs     []int32 `json:"participants"`
	From        string  `json:"from"` // 2022-01-02T15:00
	To          string  `json:"to"`   // 2022-01-02T16:00
}

type MeetingInfoResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	From        string `json:"from"`
	To          string `json:"to"`
}

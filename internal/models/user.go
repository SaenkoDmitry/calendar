package models

type CreateUserReq struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Zone       string `json:"zone"`
}

type ChangeZoneReq struct {
	Zone string `json:"zone" validate:"required"`
}

type GetMeetingsByUserAndTimeRequest struct {
	From string `json:"from"` // 2022-01-02T15:00
	To   string `json:"to"`   // 2022-01-02T16:00
}

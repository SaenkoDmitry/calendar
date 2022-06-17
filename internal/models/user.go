package models

type CreateUserReq struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
}

type ChangeZoneReq struct {
	Zone string `json:"zone" validate:"required"`
}

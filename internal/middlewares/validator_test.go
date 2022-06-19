package middlewares

import (
	"calendar/internal/models"
	"testing"

	"github.com/go-playground/validator"
)

func TestCustomValidator_Validate(t *testing.T) {
	type args interface{}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// CreateMeetingReq validation
		{
			name: "failed_without_all_required_fields_for_meeting",
			args: models.CreateMeetingReq{
				Name: "test_meeting",
			},
			wantErr: true,
		},
		{
			name: "failed_parse_date_for_meeting",
			args: models.CreateMeetingReq{
				Name:    "test_meeting",
				AdminID: 1,
				From:    "2022-01-02T15:00",
				To:      "it's_not_a_date",
			},
			wantErr: true,
		},
		{
			name: "failed_parse_date_for_meeting",
			args: models.CreateMeetingReq{
				Name:    "test_meeting",
				AdminID: 1,
				From:    "it's_not_a_date",
				To:      "2022-01-02T15:00",
			},
			wantErr: true,
		},
		{
			name: "successful_validation_for_meeting",
			args: models.CreateMeetingReq{
				Name:    "test_meeting",
				AdminID: 1,
				UserIDs: []int32{4, 5},
				From:    "2022-01-02T15:00",
				To:      "2022-01-02T15:00",
			},
			wantErr: false,
		},

		// CreateUserReq validation
		{
			name: "invalid_email_validation_for_user",
			args: models.CreateUserReq{
				FirstName:  "first",
				SecondName: "second",
				Email:      "test",
				Zone:       "Europe/Moscow",
			},
			wantErr: true,
		},
		{
			name: "missing_second_name_validation_for_user",
			args: models.CreateUserReq{
				FirstName: "first",
				Email:     "test@test.ru",
				Zone:      "Europe/Moscow",
			},
			wantErr: true,
		},
		{
			name: "successful_validation_for_user",
			args: models.CreateUserReq{
				FirstName:  "first",
				SecondName: "second",
				Email:      "test@test.ru",
				Zone:       "Europe/Moscow",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CustomValidator{Validator: validator.New()}
			if err := cv.Validate(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

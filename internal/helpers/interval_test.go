package helpers

import (
	"calendar/internal/constants"
	"testing"
)

func TestValidateRepeatInterval(t *testing.T) {
	type args struct {
		repeat string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "invalid_repeat_interval",
			args: args{
				repeat: "test",
			},
			want: false,
		},
		{
			name: "invalid_repeat_interval_days",
			args: args{
				repeat: constants.Days,
			},
			want: true,
		},
		{
			name: "valid_repeat_interval_weeks",
			args: args{
				repeat: constants.Weeks,
			},
			want: true,
		},
		{
			name: "valid_repeat_interval_months",
			args: args{
				repeat: constants.Months,
			},
			want: true,
		},
		{
			name: "valid_repeat_interval_years",
			args: args{
				repeat: constants.Years,
			},
			want: true,
		},
		{
			name: "valid_repeat_interval_weekdays",
			args: args{
				repeat: constants.Weekdays,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRepeatInterval(tt.args.repeat); got != tt.want {
				t.Errorf("ValidateRepeatInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

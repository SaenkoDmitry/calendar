package helpers

import (
	"calendar/internal/constants"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "failed_status_validation",
			args: args{
				status: "test",
			},
			want: false,
		},
		{
			name: "successful_requested_status_validation",
			args: args{
				status: constants.Requested,
			},
			want: true,
		},
		{
			name: "successful_approved_status_validation",
			args: args{
				status: constants.Approved,
			},
			want: true,
		},
		{
			name: "successful_declined_status_validation",
			args: args{
				status: constants.Declined,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidStatus(tt.args.status); got != tt.want {
				t.Errorf("IsValidStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

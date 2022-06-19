package helpers

import (
	"calendar/internal/constants"
	"reflect"
	"testing"
	"time"
)

func TestMergeDateAndTimeFromServer(t *testing.T) {
	type args struct {
		date *time.Time
		t    *time.Time
	}
	firstDate, _ := time.Parse(constants.DateTimeFormat, "2022-01-02T00:00")
	firstTime, _ := time.Parse(constants.DateTimeFormat, "2000-01-01T21:00")
	mergedDateTime, _ := time.ParseInLocation(constants.DateTimeFormat, "2022-01-02T21:00", constants.ServerTimeZone)
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "merged_date_and_time",
			args: args{
				date: &firstDate,
				t:    &firstTime,
			},
			want: mergedDateTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeDateAndTimeFromServer(tt.args.date, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeDateAndTimeFromServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

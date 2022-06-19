package helpers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestChooseZone(t *testing.T) {
	moscowLoc, _ := time.LoadLocation("Europe/Moscow")

	type args struct {
		c    echo.Context
		zone string
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	tests := []struct {
		name    string
		args    args
		want    *time.Location
		wantErr bool
	}{
		{
			name: "invalid_zone_validation",
			args: args{
				c:    c,
				zone: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "successful_zone_validation",
			args: args{
				c:    c,
				zone: "Europe/Moscow",
			},
			want:    moscowLoc,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChooseZone(tt.args.c, tt.args.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChooseZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChooseZone() got = %v, want %v", got, tt.want)
			}
		})
	}
}

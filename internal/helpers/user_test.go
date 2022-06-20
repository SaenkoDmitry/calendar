package helpers

import (
	"reflect"
	"testing"
)

func TestCheckMissingUserIDs(t *testing.T) {
	type args struct {
		passed  []int32
		existed []int32
	}
	type want struct {
		res0 []int32
		res1 bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "passed_the_same_as_existed",
			args: args{
				passed:  []int32{1, 2, 3, 4, 5},
				existed: []int32{1, 2, 3, 4, 5},
			},
			want: want{
				res0: []int32{},
				res1: true,
			},
		},
		{
			name: "passed_bigger_than_existed",
			args: args{
				passed:  []int32{1, 2, 3, 4, 5, 6},
				existed: []int32{1, 2, 3, 4, 5},
			},
			want: want{
				res0: []int32{6},
				res1: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ValidateExistenceOfUsers(tt.args.passed, tt.args.existed)
			if !reflect.DeepEqual(got, tt.want.res0) {
				t.Errorf("ValidateExistenceOfUsers() got = %v, want %v", got, tt.want.res0)
			}
			if got1 != tt.want.res1 {
				t.Errorf("ValidateExistenceOfUsers() got = %v, want %v", got1, tt.want.res1)
			}
		})
	}
}

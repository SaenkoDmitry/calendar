package helpers

import (
	"calendar/internal/constants"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/jackc/pgconn"
)

func IsEmailDuplicated(err *pgconn.PgError) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Message, "duplicate key")
}

func GetUser(c echo.Context, paramName string) (int32, error) {
	temp, err := strconv.ParseInt(c.Param(paramName), 10, 32)
	if err != nil || temp <= 0 {
		return 0, WrapError(c, http.StatusBadRequest, constants.InvalidUserID)
	}
	return int32(temp), nil
}

// ValidateExistenceOfUsers extract users from first that not exists in second
func ValidateExistenceOfUsers(first, second []int32) ([]int32, bool) {
	res := make([]int32, 0)
	m1 := make(map[int32]struct{})
	m2 := make(map[int32]struct{})
	for i := range first {
		m1[first[i]] = struct{}{}
	}
	for i := range second {
		m2[second[i]] = struct{}{}
	}
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			res = append(res, k)
		}
	}
	if len(res) > 0 {
		return res, false
	}
	return res, true
}

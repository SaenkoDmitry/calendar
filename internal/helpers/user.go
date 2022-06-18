package helpers

import (
	"strings"

	"github.com/jackc/pgconn"
)

func IsEmailDuplicated(err *pgconn.PgError) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Message, "duplicate key")
}

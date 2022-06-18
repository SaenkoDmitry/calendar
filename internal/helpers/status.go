package helpers

import "calendar/internal/constants"

func IsValidStatus(status string) bool {
	for i := range constants.ValidStatuses {
		if status == constants.ValidStatuses[i] {
			return true
		}
	}
	return false
}

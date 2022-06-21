package helpers

import (
	"calendar/internal/constants"
	"calendar/internal/models"
	"time"
)

func ValidateRepeatInterval(repeat string) bool {
	for i := range constants.ValidRepeatIntervals {
		if constants.ValidRepeatIntervals[i] == repeat {
			return true
		}
	}
	return false
}

func CalcVirtualMeetingsInsideInterval(meet *models.VirtualMeetingInfo, from time.Time, to time.Time) []*models.VirtualMeetingInfo {
	switch meet.Repeat {
	case constants.Days:
		return handleDays(meet, from, to)
	case constants.Weeks:
		return handleWeeks(meet, from, to)
	case constants.Months:
		return handleMonths(meet, from, to)
	case constants.Years:
		return handleYears(meet, from, to)
	case constants.Weekdays:
		return handleWeekdays(meet, from, to)
	}
	return []*models.VirtualMeetingInfo{}
}

func handleWeekdays(meet *models.VirtualMeetingInfo, from time.Time, to time.Time) []*models.VirtualMeetingInfo {
	res := make([]*models.VirtualMeetingInfo, 0)
	return res
}

func handleYears(meet *models.VirtualMeetingInfo, from time.Time, to time.Time) []*models.VirtualMeetingInfo {
	res := make([]*models.VirtualMeetingInfo, 0)
	return res
}

func handleMonths(meet *models.VirtualMeetingInfo, from time.Time, to time.Time) []*models.VirtualMeetingInfo {
	res := make([]*models.VirtualMeetingInfo, 0)
	return res
}

func handleWeeks(meet *models.VirtualMeetingInfo, from time.Time, to time.Time) []*models.VirtualMeetingInfo {
	res := make([]*models.VirtualMeetingInfo, 0)
	return res
}

func handleDays(meet *models.VirtualMeetingInfo, from time.Time, to time.Time) []*models.VirtualMeetingInfo {
	res := make([]*models.VirtualMeetingInfo, 0)
	return res
}

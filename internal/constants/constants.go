package constants

import "time"

const (
	EmailAlreadyRegistered = "email_already_registered"
	UndefinedDB            = "undefined_db_error"
	NothingUpdated         = "nothing_updated"

	InvalidUserID   = "invalid_user_id"
	UserIDNotExists = "user_id_not_exists"
	InvalidUserList = "invalid_user_list"

	InvalidOrEmptyStatus = "invalid_or_empty_status"

	InvalidMeetingID                     = "invalid_meeting_id"
	MeetingIDNotExists                   = "meeting_id_not_exists"
	CannotCreateMeeting                  = "cannot_create_meeting"
	TooManyUsersForMeeting               = "too_many_users_for_meeting"
	InvalidRepeatIntervals               = "invalid_repeat_intervals"
	NotFoundOptimalMeetingForTheInterval = "not_found_optimal_meeting_for_the_interval"

	NotValidTimeZone      = "invalid_time_zone"
	InvalidFromDate       = "invalid_from_date"
	InvalidToDate         = "invalid_to_date"
	FromEarlierThanToDate = "from_earlier_than_to_date"
)

const (
	DateTimeFormat       = "2006-01-02T15:04"
	DateFormat           = "2006-01-02"
	TimeFormat           = "15:04"
	PrettyDateTimeFormat = "Mon, 02 Jan 2006 15:04"
)

const (
	UserIDConstraintDBErr = `insert or update on table "user_meetings" violates foreign key constraint "fk_user_id"`
)

var ValidStatuses = []string{Requested, Approved, Declined}

const (
	Requested = "requested"
	Approved  = "approved"
	Declined  = "declined"
)

const (
	Days     = "days"
	Weeks    = "weeks"
	Months   = "months"
	Years    = "years"
	Weekdays = "weekdays"
)

var ValidRepeatIntervals = []string{Days, Weeks, Months, Years, Weekdays}

const (
	PostgresDBService = "postgres_db"
)

var (
	ServerTimeZone *time.Location
)

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	ServerTimeZone = loc
}

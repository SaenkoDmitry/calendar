package constants

const (
	EmailAlreadyRegistered = "email_already_registered"
	UndefinedDB            = "undefined_db_error"

	UserIDNotExists     = "user_id_not_exists"
	MeetingIDNotExists  = "meeting_id_not_exists"
	CannotInsertMeeting = "cannot_insert_meeting"

	NotValidTimeZone      = "not_valid_time_zone"
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

package constants

import "time"

const (
	EmailAlreadyRegistered = "email_already_registered"
	UndefinedDB            = "undefined_db_error"

	InvalidUserID   = "invalid_user_id"
	UserIDNotExists = "user_id_not_exists"

	EmptyStatus   = "not_passed_status"
	InvalidStatus = "invalid_status"

	InvalidMeetingID                     = "invalid_meeting_id"
	MeetingIDNotExists                   = "meeting_id_not_exists"
	CannotCreateMeeting                  = "cannot_create_meeting"
	TooManyUsersForMeeting               = "too_many_users_for_meeting"
	UserNotInvitedOnTheMeeting           = "user_not_invited_on_the_meeting"
	MeetingCanceledOrFinished            = "meeting_canceled_or_finished"
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

var ValidStatuses = []string{Requested, Approved, Declined, Finished, Canceled}

const (
	Requested = "requested"
	Approved  = "approved"
	Declined  = "declined"
	Finished  = "finished"
	Canceled  = "canceled"
)

const (
	PostgresDBService = "postgres_db"
)

var (
	ServerTimeZone *time.Location
)

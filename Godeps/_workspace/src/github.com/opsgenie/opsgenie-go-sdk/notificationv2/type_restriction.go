package notificationv2

const (
	// List of time restrictions within alerts will be sent.
	TimeOfDayTimeRestriction           = "time-of-day"
	WeekendAndTimeOfDayTimeRestriction = "weekday-and-time-of-day"
)

// TypeRestriction is a type of restriction, within alerts will be sent.
type TypeRestriction string

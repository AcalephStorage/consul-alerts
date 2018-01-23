package schedulev2

const (
	// The list of week days. These strings are used for generate time restrictions.
	Monday    Day = "monday"
	Tuesday   Day = "tuesday"
	Wednesday Day = "wednesday"
	Thursday  Day = "thursday"
	Friday    Day = "friday"
	Saturday  Day = "saturday"
	Sunday    Day = "sunday"
)

// Day is the text representation of day name of week.
type Day string

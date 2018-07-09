package schedulev2

type TimeRestriction struct {
	Type         RestrictionType `json:"type"`
	Restriction  Restriction     `json:"restriction"`
	Restrictions []Restriction   `json:"restrictions"`
}

type Restriction struct {
	StartHour int `json:"startHour,omitempty"`
	StartMin  int `json:"startMin,omitempty"`
	StartDay  Day `json:"startDay,omitempty"`
	EndHour   int `json:"endHour,omitempty"`
	EndMin    int `json:"endMin,omitempty"`
	EndDay    Day `json:"endDay,omitempty"`
}

const (
	DayRestrictionType     RestrictionType = "time-of-day"
	WeekDayRestrictionType RestrictionType = "weekday-and-time-of-day"
)

type RestrictionType string

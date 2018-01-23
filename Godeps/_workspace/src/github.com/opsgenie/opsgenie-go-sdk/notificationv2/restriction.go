package notificationv2

// TimeRestriction is used to limit notification rules to certain day and time of the week, using multiple start and
// end times for each day of the week.
type TimeRestriction struct {
	Type         TypeRestriction `json:"type,omitempty"`
	Restrictions []Restriction   `json:"restrictions,omitempty"`
	Restriction  Restriction     `json:"restriction, omitempty"`
}

// Restriction defines start and end times.
type Restriction struct {
	StartDay  Day  `json:"startDay,omitempty"`
	EndDay    Day  `json:"endDay,omitempty"`
	StartHour int `json:"startHour,omitempty"`
	EndHour   int `json:"endHour,omitempty"`
	StartMin  int `json:"startMin,omitempty"`
	EndMin    int `json:"endMin,omitempty"`
}

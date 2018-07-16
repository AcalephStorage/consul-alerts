package schedulev2

type Rotation struct {
	Name            string          `json:"name,omitempty"`
	StartDate       string          `json:"startDate,omitempty"`
	EndDate         string          `json:"endDate,omitempty"`
	Type            Type            `json:"type,omitempty"`
	Length          int             `json:"length,omitempty"`
	Participants    []Participant   `json:"participants,omitempty"`
	TimeRestriction TimeRestriction `json:"timeRestriction,omitempty"`
}

const (
	DailyRotation  Type = "daily"
	WeeklyRotation Type = "weekly"
	HourlyRotation Type = "hourly"
)

type Type string

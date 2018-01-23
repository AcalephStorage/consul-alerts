package schedulev2


type ScheduleOverride struct {
	Parent           Parent             `json:"_parent,omitempty"`
	Alias            string             `json:"alias,omitempty"`
	User             User               `json:"user,omitempty"`
	StartDate        string             `json:"startDate,omitempty"`
	EndDate          string     	    `json:"endDate,omitempty"`
	Rotations		 []Rotation			`json:"rotations,omitempty"`
}


type Parent struct {
	ID			string 			`json:"id,omitempty"`
	Name        string          `json:"name,omitempty"`
	Enabled     bool            `json:"enabled,omitempty"`
}

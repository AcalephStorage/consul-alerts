package schedulev2


type Schedule struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Description      string     	    `json:"description,omitempty"`
	Timezone		 string 			`json:"timezone,omitempty"`
	Enabled          bool               `json:"enabled,omitempty"`
	OwnerTeam        OwnerTeam          `json:"ownerTeam,omitempty"`
	Rotations		 []Rotation			`json:"rotations,omitempty"`
}

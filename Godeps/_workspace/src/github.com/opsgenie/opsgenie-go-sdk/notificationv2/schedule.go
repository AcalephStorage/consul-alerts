package notificationv2

// Schedule defines name and id of schedule. The field "type" is mandatory and should be set as "schedule".
type Schedule struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

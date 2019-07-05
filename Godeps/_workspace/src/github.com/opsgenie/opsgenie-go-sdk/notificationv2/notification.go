package notificationv2

// Notification describes alert data, which is contained in a response.
type Notification struct {
	ID               string             `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	ActionType       ActionType         `json:"actionType,omitempty"`
	NotificationTime []NotificationTime `json:"notificationTime,omitempty"`
	Order            int                `json:"order,omitempty"`
	Steps            []Step             `json:"steps,omitempty"`
	Schedules        []Schedule         `json:"schedules,omitempty"`
	Criteria         Criteria           `json:"criteria,omitempty"`
	Enabled          bool               `json:"enabled,omitempty"`
	Repeat           Repeat             `json:"repeat,omitempty"`
	TimeRestriction  TimeRestriction    `json:"timeRestriction,omitempty"`
}

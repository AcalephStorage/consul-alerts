package userv2

// ListUserSchedulesResponse is a response with list of user schedules.
type ListUserSchedulesResponse struct {
	Schedules []Schedule `json:"data,omitempty"`
	ResponseMeta
}

// Schedule contains data of schedule.
type Schedule struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

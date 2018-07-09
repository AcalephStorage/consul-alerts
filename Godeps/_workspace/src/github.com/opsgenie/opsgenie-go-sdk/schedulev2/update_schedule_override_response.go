package schedulev2

// UpdateScheduleOverrideResponse is a response of get alert action.
type UpdateScheduleOverrideResponse struct {
	ResponseMeta
	ScheduleOverride ScheduleOverride `json:"data"`
}

package schedulev2

// GetScheduleOverrideResponse is a response of get alert action.
type GetScheduleOverrideResponse struct {
	ResponseMeta
	ScheduleOverride ScheduleOverride `json:"data"`
}

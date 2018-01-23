package schedulev2

// ListScheduleOverrideResponse is a response of get alert action.
type ListScheduleOverrideResponse struct {
	ResponseMeta
	ScheduleOverrides []ScheduleOverride `json:"data"`
}


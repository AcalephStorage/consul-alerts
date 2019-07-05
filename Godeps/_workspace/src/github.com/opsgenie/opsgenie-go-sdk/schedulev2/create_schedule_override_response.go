package schedulev2

// CreateScheduleOverrideResponse is a response of creating alert action.
type CreateScheduleOverrideResponse struct {
	ResponseMeta
	ScheduleOverride ScheduleOverride `json:"data"`
}

package schedulev2

// CreateScheduleResponse is a response of creating alert action.
type CreateScheduleResponse struct {
	ResponseMeta
	Schedule Schedule `json:"data"`
}

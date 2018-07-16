package schedulev2

// CreateScheduleRotationResponse is a response of creating alert action.
type CreateScheduleRotationResponse struct {
	ResponseMeta
	Schedule Schedule `json:"data"`
}

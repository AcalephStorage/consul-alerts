package schedulev2

// GetScheduleRotationResponse is a response of get alert action.
type GetScheduleRotationResponse struct {
	ResponseMeta
	Schedule Schedule `json:"data"`
}

package schedulev2

// UpdateScheduleRotationResponse is a response of get alert action.
type UpdateScheduleRotationResponse struct {
	ResponseMeta
	Schedule Schedule `json:"data"`
}

package schedulev2

// ListScheduleRotationResponse is a response of get alert action.
type ListScheduleRotationResponse struct {
	ResponseMeta
	Schedule []Schedule `json:"data"`
}

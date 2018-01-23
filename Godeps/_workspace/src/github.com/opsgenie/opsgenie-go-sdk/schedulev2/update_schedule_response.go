package schedulev2

// UpdateScheduleResponse is a response of get alert action.
type UpdateScheduleResponse struct {
	ResponseMeta
	Schedule Schedule `json:"data"`
}

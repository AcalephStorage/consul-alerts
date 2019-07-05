package schedulev2

// ListScheduleResponse is a response of get alert action.
type ListScheduleResponse struct {
	ResponseMeta
	Schedule []Schedule `json:"data"`
}

package schedulev2

// GetScheduleResponse is a response of get alert action.
type GetScheduleResponse struct {
	ResponseMeta
	Schedule Schedule `json:"data"`
}


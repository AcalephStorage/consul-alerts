package alertsv2

type ResponseMeta struct {
	RequestID      string
	ResponseTime   float32
	RateLimitState string
}

func (rm *ResponseMeta) SetRequestID(requestID string) {
	rm.RequestID = requestID
}

func (rm *ResponseMeta) SetResponseTime(responseTime float32) {
	rm.ResponseTime = responseTime
}

func (rm *ResponseMeta) SetRateLimitState(state string) {
	rm.RateLimitState = state
}

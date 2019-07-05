package userv2

// ResponseMeta is a struct of meta data.
type ResponseMeta struct {
	RequestID      string
	ResponseTime   float32
	RateLimitState string
}

// SetRequestID sets identifier of request.
func (rm *ResponseMeta) SetRequestID(requestID string) {
	rm.RequestID = requestID
}

// SetResponseTime sets response time.
func (rm *ResponseMeta) SetResponseTime(responseTime float32) {
	rm.ResponseTime = responseTime
}

// SetRateLimitState sets state of rate limit.
func (rm *ResponseMeta) SetRateLimitState(state string) {
	rm.RateLimitState = state
}

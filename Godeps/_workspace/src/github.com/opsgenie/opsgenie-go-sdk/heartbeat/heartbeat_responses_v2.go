package heartbeat

import "time"

type HeartbeatResponseV2 struct {
	Data      HeartbeatData `json:"data"`
	Took      float32 `json:"took"`
	RequestId string   `json:"requestId"`
}

func (rm *HeartbeatResponseV2) SetRequestID(requestID string) {
	rm.RequestId = requestID
}

func (rm *HeartbeatResponseV2) SetResponseTime(responseTime float32) {
	rm.Took = responseTime
}

func (rm *HeartbeatResponseV2) SetRateLimitState(state string) {
}

type HeartbeatData struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	LastHeartbeat time.Time `json:"lastPingTime,omitempty"`
	Enabled       bool   `json:"enabled"`
	Interval      int    `json:"interval"`
	IntervalUnit  string `json:"intervalUnit"`
	Expired       bool   `json:"expired"`
}

type HeartbeatMetaResponseV2 struct {
	Code      int `json:"code"`
	Data      HeartbeatMetaData `json:"data"`
	Took      float32 `json:"took"`
	RequestId string   `json:"requestId"`
}

func (rm *HeartbeatMetaResponseV2) SetRequestID(requestID string) {
	rm.RequestId = requestID
}

func (rm *HeartbeatMetaResponseV2) SetResponseTime(responseTime float32) {
	rm.Took = responseTime
}

func (rm *HeartbeatMetaResponseV2) SetRateLimitState(state string) {
}

type HeartbeatMetaData struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Expired bool   `json:"expired"`
}

type PingHeartbeatResponse struct {
	Result    string `json:"result"`
	Took      float32 `json:"took"`
	RequestId string   `json:"requestId"`
}

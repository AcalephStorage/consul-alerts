package alertsv2

type Report struct {
	AckTime        int32  `json:"ackTime,omitempty"`
	CloseTime      int32  `json:"closeTime,omitempty"`
	AcknowledgedBy string `json:"acknowledgedBy,omitempty"`
	ClosedBy       string `json:"closedBy,omitempty"`
}

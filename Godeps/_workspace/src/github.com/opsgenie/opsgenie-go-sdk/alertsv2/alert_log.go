package alertsv2

import "time"

type AlertLog struct {
	Log       string    `json:"log,omitempty"`
	Type      string    `json:"type,omitempty"`
	Owner     string    `json:"owner,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Offset    string    `json:"offset,omitempty"`
}

type ListAlertLogsResponse struct {
	ResponseMeta
	AlertLogs []AlertLog `json:"data"`
}

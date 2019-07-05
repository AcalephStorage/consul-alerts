package alertsv2

import "time"

type AlertNote struct {
	Note      string    `json:"note,omitempty"`
	Owner     string    `json:"owner,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Offset    string    `json:"offset,omitempty"`
}

type ListAlertNotesResponse struct {
	ResponseMeta
	AlertNotes []AlertNote `json:"data"`
}

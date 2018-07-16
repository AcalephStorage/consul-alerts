package alertsv2

import "time"

type AlertRecipient struct {
	User      User      `json:"user,omitempty"`
	State     string    `json:"state,omitempty"`
	Method    string    `json:"method,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type ListAlertRecipientsResponse struct {
	ResponseMeta
	Recipients []AlertRecipient `json:"data"`
}

package alertsv2

import "time"

type RequestStatus struct {
	IsSuccess     bool `json:"isSuccess,omitempty"`
	Action        string `json:"action,omitempty"`
	ProcessedAt   time.Time `json:"processedAt,omitempty"`
	IntegrationId string `json:"integrationId,omitempty"`
	Status        string `json:"status,omitempty"`
	AlertID       string `json:"alertId,omitempty"`
	Alias         string `json:"alias,omitempty"`
}

type GetAsyncRequestStatusResponse struct {
	ResponseMeta
	Status       RequestStatus `json:"data"`
}

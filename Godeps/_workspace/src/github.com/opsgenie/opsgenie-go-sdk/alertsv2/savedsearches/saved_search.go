package savedsearches

import (
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/alertsv2"
	"time"
)

type SavedSearch struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	CreatedAt   time.Time           `json:"createdAt,omitempty"`
	UpdatedAt   time.Time           `json:"updatedAt,omitempty"`
	Teams       []alertsv2.TeamMeta `json:"teams,omitempty"`
	Description string              `json:"description,omitempty"`
	Query       string              `json:"query,omitempty"`
}

type SavedSearchResponse struct {
	alertsv2.ResponseMeta
	SavedSearch SavedSearch `json:"data"`
}

type GetSavedSearchResponse struct {
	alertsv2.ResponseMeta
	SavedSearch SavedSearch `json:"data"`
}

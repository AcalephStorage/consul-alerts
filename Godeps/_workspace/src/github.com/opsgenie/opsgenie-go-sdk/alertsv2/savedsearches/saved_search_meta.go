package savedsearches

import "github.com/opsgenie/opsgenie-go-sdk/alertsv2"

type SavedSearchMeta struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CreateSavedSearchResponse struct {
	alertsv2.ResponseMeta
	SavedSearch SavedSearchMeta `json:"data"`
}

type ListSavedSearchResponse struct {
	alertsv2.ResponseMeta
	SavedSearches []SavedSearchMeta `json:"data"`
}

type UpdateSavedSearchResponse struct {
	alertsv2.ResponseMeta
	SavedSearch SavedSearchMeta `json:"data"`
}

type DeleteSavedSearchResponse struct {
	alertsv2.ResponseMeta
	SavedSearch SavedSearchMeta `json:"data"`
}

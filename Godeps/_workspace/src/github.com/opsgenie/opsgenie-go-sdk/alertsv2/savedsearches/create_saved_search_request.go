package savedsearches

import (
	"net/url"
	"github.com/opsgenie/opsgenie-go-sdk/alertsv2"
)

type CreateSavedSearchRequest struct {
	Name        string `json:"name,omitempty"`
	Query       string `json:"query,omitempty"`
	Owner       alertsv2.User `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Teams       []alertsv2.Team `json:"teams,omitempty"`
	ApiKey      string `json:"apiKey,omitempty"`
}

func (r *CreateSavedSearchRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *CreateSavedSearchRequest) GenerateUrl() (string, url.Values, error) {
	return "/v2/alerts/saved-searches", nil, nil
}

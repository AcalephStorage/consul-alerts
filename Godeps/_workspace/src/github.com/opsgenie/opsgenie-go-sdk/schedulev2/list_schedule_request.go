package schedulev2

import (
	"net/url"
)

// ListScheduleRequest is a struct of request to crate new schedule.
type ListScheduleRequest struct {
	ApiKey           string
	Expand			 string
}

// GetApiKey returns api key.
func (r *ListScheduleRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *ListScheduleRequest) GenerateUrl() (string, url.Values, error) {
	params := url.Values{}

	if r.Expand != ""{
		params.Add("expand", r.Expand)
	}

	return "/v2/schedules/", params, nil
}

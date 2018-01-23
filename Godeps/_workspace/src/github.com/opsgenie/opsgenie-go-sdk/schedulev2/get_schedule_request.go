package schedulev2

import (
	"net/url"
)

// GetScheduleRequest is a struct of request to crate new schedule.
type GetScheduleRequest struct {
	*Identifier
	ApiKey           string
}

// GetApiKey returns api key.
func (r *GetScheduleRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *GetScheduleRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	return baseUrl, params, nil
}

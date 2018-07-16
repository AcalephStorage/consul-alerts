package schedulev2

import (
	"net/url"
)

// DeleteScheduleRequest is a struct of request to delete schedule.
type DeleteScheduleRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *DeleteScheduleRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *DeleteScheduleRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	return baseUrl, params, nil
}

package schedulev2

import (
	"net/url"
)

// ListScheduleOverrideRequest is a struct of request to crate new schedule.
type ListScheduleOverrideRequest struct {
	*ScheduleIdentifier
	ApiKey           string
}

// GetApiKey returns api key.
func (r *ListScheduleOverrideRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *ListScheduleOverrideRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	baseUrl += "/rotations"
	return baseUrl, params, nil
}

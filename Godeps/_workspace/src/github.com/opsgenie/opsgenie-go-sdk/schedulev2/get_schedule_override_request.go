package schedulev2

import (
	"errors"
	"net/url"
)

// GetScheduleOverrideRequest is a struct of request to crate new schedule.
type GetScheduleOverrideRequest struct {
	*ScheduleIdentifier
	ApiKey string
	Alias  string
}

// GetApiKey returns api key.
func (r *GetScheduleOverrideRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *GetScheduleOverrideRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if r.Alias == "" {
		return "", nil, errors.New("schedule alias should be provided for get action")
	}
	baseUrl += "/overrides"

	return baseUrl + "/" + r.Alias, params, nil
}

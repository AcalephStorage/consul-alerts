package schedulev2

import (
	"net/url"
	"errors"
)

// DeleteScheduleOverrideRequest is a struct of request to delete schedule.
type DeleteScheduleOverrideRequest struct {
	*ScheduleIdentifier
	Alias		     string
	ApiKey           string
}

// GetApiKey returns api key.
func (r *DeleteScheduleOverrideRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *DeleteScheduleOverrideRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if r.Alias == "" {
		return "", nil, errors.New("alias should be provided for delete action")
	}
	baseUrl += "/overrides"

	return baseUrl + "/" + r.Alias, params, nil
}

package schedulev2

import (
	"errors"
	"net/url"
)

// GetScheduleRotationRequest is a struct of request to crate new schedule.
type GetScheduleRotationRequest struct {
	*ScheduleIdentifier
	ApiKey string
	ID     string
}

// GetApiKey returns api key.
func (r *GetScheduleRotationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *GetScheduleRotationRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if r.ID == "" {
		return "", nil, errors.New("schedule ID should be provided for get action")
	}
	baseUrl += "/rotations"

	return baseUrl + "/" + r.ID, params, nil
}

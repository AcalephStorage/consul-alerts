package schedulev2

import (
	"net/url"
	"errors"
)

// DeleteScheduleRotationRequest is a struct of request to delete schedule.
type DeleteScheduleRotationRequest struct {
	*ScheduleIdentifier
	ID				 string
	ApiKey           string
}

// GetApiKey returns api key.
func (r *DeleteScheduleRotationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *DeleteScheduleRotationRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if r.ID == "" {
		return "", nil, errors.New("schedule ID should be provided for delete action")
	}
	baseUrl += "/rotations"

	return baseUrl + "/" + r.ID, params, nil
}

package schedulev2

import (
	"net/url"
	"errors"
)

// UpdateScheduleOverrideRequest is a struct of request to update existing schedule.
type UpdateScheduleOverrideRequest struct {
	*ScheduleIdentifier
	ApiKey           string
	Alias            string
	User			 User				`json:"user"`
	StartDate		 string				`json:"startDate"`
	EndDate		 	 string				`json:"endDate"`
	Rotations		 []Rotation			`json:"rotations"`
}

// GetApiKey returns api key.
func (r *UpdateScheduleOverrideRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *UpdateScheduleOverrideRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if r.Alias == "" {
		return "", nil, errors.New("alias should be provided for update action")
	}

	if r.StartDate == "" {
		return "", nil, errors.New("start date should be provided for update action")
	}

	if r.EndDate == "" {
		return "", nil, errors.New("end date should be provided for update action")
	}

	baseUrl += "/overrides"
	return  baseUrl + "/" + r.Alias , params, nil
}

package schedulev2

import (
	"errors"
	"net/url"
)

// CreateScheduleOverrideRequest is a struct of request to crate new schedule.
type CreateScheduleOverrideRequest struct {
	*ScheduleIdentifier
	ApiKey    string
	Alias     string     `json:"alias"`
	User      User       `json:"user"`
	StartDate string     `json:"startDate"`
	EndDate   string     `json:"endDate"`
	Rotations []Rotation `json:"rotations"`
}

// GetApiKey returns api key.
func (r *CreateScheduleOverrideRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *CreateScheduleOverrideRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if r.Alias == "" {
		return "", nil, errors.New("alias should be provided for create action")
	}

	if r.User.Type == UserUserType && (r.User.ID == "" && r.User.Username == "") {
		return "", nil, errors.New("alias should be provided for create action")
	}

	if r.StartDate == "" {
		return "", nil, errors.New("StartDate should be provided for create action")
	}

	if r.EndDate == "" {
		return "", nil, errors.New("EndDate should be provided for create action")
	}

	baseUrl += "/overrides"
	return baseUrl, params, nil
}

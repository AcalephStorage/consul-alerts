package schedulev2

import (
	"errors"
	"net/url"
)

// CreateScheduleRotationRequest is a struct of request to crate new schedule.
type CreateScheduleRotationRequest struct {
	*ScheduleIdentifier
	ApiKey          string
	Name            string          `json:"name"`
	StartDate       string          `json:"startDate"`
	EndDate         string          `json:"endDate"`
	Type            Type            `json:"type"`
	Length          int             `json:"length"`
	Participants    []Participant   `json:"participants"`
	TimeRestriction TimeRestriction `json:"timeRestriction"`
}

// GetApiKey returns api key.
func (r *CreateScheduleRotationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *CreateScheduleRotationRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if r.StartDate == "" {
		return "", nil, errors.New("StartDate should be provided for create action")
	}

	if r.Type == "" {
		return "", nil, errors.New("Type should be provided for create action")
	}

	if len(r.Participants) < 1 {
		return "", nil, errors.New("At least one Participants should be provided for create action")
	}
	baseUrl += "/rotations"

	return baseUrl, params, nil
}

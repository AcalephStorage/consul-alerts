package schedulev2

import (
	"net/url"
	"errors"
)

// UpdateScheduleRotationRequest is a struct of request to update existing schedule.
type UpdateScheduleRotationRequest struct {
	*ScheduleIdentifier
	ApiKey           string
	ID				 string
	Name             string             `json:"name"`
	StartDate		 string				`json:"startDate"`
	EndDate		 	 string				`json:"endDate"`
	Type		 	 Type    			`json:"type"`
	Length		 	 int				`json:"length"`
	Participants     []Participant		`json:"participants"`
	TimeRestriction  TimeRestriction    `json:"timeRestriction"`
}

// GetApiKey returns api key.
func (r *UpdateScheduleRotationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *UpdateScheduleRotationRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.ScheduleIdentifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if r.ID == "" {
		return "", nil, errors.New("schedule ID should be provided for update action")
	}

	baseUrl += "/rotations"
	return  baseUrl + "/" + r.ID , params, nil
}

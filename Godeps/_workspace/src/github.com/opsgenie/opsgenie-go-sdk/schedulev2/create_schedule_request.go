package schedulev2

import (
	"net/url"
	"errors"
)

// CreateScheduleRequest is a struct of request to crate new schedule.
type CreateScheduleRequest struct {
	ApiKey           string
	Name             string             `json:"name"`
	Description		 string            	`json:"description"`
	Timezone		 string				`json:"timezone"`
	Enabled          bool               `json:"enabled"`
	OwnerTeam		 OwnerTeam			`json:"ownerTeam"`
	Rotations		 []Rotation			`json:"rotations"`
}

// GetApiKey returns api key.
func (r *CreateScheduleRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *CreateScheduleRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/schedules"

	if r.Name == "" {
		return "", nil, errors.New("Name should be provided for create action")
	}

	if len(r.Rotations) < 1 {
		return "", nil, errors.New("At least one roation should be provided for create action")
	}

	return baseUrl, nil, nil
}

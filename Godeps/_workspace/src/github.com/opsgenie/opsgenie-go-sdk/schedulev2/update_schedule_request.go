package schedulev2

import (
	"errors"
	"net/url"
)

// UpdateScheduleRequest is a struct of request to update existing schedule.
type UpdateScheduleRequest struct {
	*Identifier
	ApiKey      string
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Timezone    string     `json:"timezone"`
	Enabled     bool       `json:"enabled"`
	OwnerTeam   OwnerTeam  `json:"ownerTeam"`
	Rotations   []Rotation `json:"rotations"`
}

// GetApiKey returns api key.
func (r *UpdateScheduleRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *UpdateScheduleRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()
	if err != nil {
		return "", nil, err
	}

	if len(r.Rotations) < 1 {
		return "", nil, errors.New("At least one roation should be provided for update action")
	}

	return baseUrl, params, nil
}

package userv2

import (
	"net/url"
)

// ListUserTeamsRequest is a request for getting list of user teams.
type ListUserTeamsRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *ListUserTeamsRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates API url using specified attributes of identifier.
func (r *ListUserTeamsRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return baseUrl + "/teams", params, err
}

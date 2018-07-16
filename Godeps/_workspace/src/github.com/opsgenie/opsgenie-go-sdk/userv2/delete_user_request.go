package userv2

import (
	"net/url"
)

// DeleteUserRequest is a request for deleting user.
type DeleteUserRequest struct {
	*Identifier
	ApiKey string
}

// GenerateUrl generates url to API endpoint.
func (r *DeleteUserRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return baseUrl, nil, err
}

// GetApiKey returns api key.
func (r *DeleteUserRequest) GetApiKey() string {
	return r.ApiKey
}

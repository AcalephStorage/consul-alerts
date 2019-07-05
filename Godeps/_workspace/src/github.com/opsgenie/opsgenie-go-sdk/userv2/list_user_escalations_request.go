package userv2

import (
	"net/url"
)

// ListUserEscalationsRequest is a request for getting user escalation list.
type ListUserEscalationsRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *ListUserEscalationsRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates API url using specified attributes of identifier.
func (r *ListUserEscalationsRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	baseUrl += "/escalations"

	return baseUrl, params, nil
}

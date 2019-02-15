package userv2

import (
	"net/url"
)

// ListUserForwardingRulesRequest is a request for getting list of forwarding rules.
type ListUserForwardingRulesRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *ListUserForwardingRulesRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates API url using specified attributes of identifier.
func (r *ListUserForwardingRulesRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return baseUrl + "/forwarding-rules", params, err
}

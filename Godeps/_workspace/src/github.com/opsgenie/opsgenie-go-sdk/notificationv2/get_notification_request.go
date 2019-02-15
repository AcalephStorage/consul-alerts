package notificationv2

import (
	"errors"
	"net/url"
)

// GetNotificationRequest is a struct of request to get notification rule.
type GetNotificationRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *GetNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *GetNotificationRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if r.Identifier.RuleID == "" {
		return "", nil, errors.New("Rule ID should be provided for get action")
	}

	baseUrl += "/" + r.Identifier.RuleID

	return baseUrl, nil, nil
}

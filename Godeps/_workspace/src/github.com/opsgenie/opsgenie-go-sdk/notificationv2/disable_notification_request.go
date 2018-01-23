package notificationv2

import (
	"net/url"
	"errors"
)

// DisableNotificationRequest is a struct of request to disable specified notification rule.
type DisableNotificationRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *DisableNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *DisableNotificationRequest) GenerateUrl() (string, url.Values, error) {

	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "" , nil, err
	}

	if r.Identifier.RuleID == "" {
		return "", nil, errors.New("Rule ID should be provided for disable action")
	}

	baseUrl += "/" + r.Identifier.RuleID + "/disable"

	return baseUrl, nil, nil
}

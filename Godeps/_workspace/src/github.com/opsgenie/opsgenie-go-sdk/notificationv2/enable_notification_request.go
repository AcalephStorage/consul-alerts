package notificationv2

import (
	"net/url"
	"errors"
)

// EnableNotificationRequest is a struct of request to enable specified notification rule.
type EnableNotificationRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *EnableNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *EnableNotificationRequest) GenerateUrl() (string, url.Values, error) {

	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "" , nil, err
	}

	if r.Identifier.RuleID == "" {
		return "", nil, errors.New("Rule ID should be provided for enable action")
	}

	baseUrl += "/" + r.Identifier.RuleID + "/enable"

	return baseUrl, nil, nil
}

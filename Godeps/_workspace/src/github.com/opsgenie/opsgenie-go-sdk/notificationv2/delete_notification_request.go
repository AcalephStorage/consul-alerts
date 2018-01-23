package notificationv2

import (
	"net/url"
	"errors"
)

// DeleteNotificationRequest is a struct of request to delete existing notification rule.
type DeleteNotificationRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *DeleteNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *DeleteNotificationRequest) GenerateUrl() (string, url.Values, error) {

	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "" , nil, err
	}

	if r.Identifier.RuleID == "" {
		return "", nil, errors.New("Rule ID should be provided for delete action")
	}

	baseUrl += "/" + r.Identifier.RuleID

	return baseUrl, nil, nil
}

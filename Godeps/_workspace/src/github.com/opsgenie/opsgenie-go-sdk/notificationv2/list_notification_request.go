package notificationv2

import (
	"net/url"
)

// ListNotificationRequest is a struct of request to get list of existing notification rules.
type ListNotificationRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *ListNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *ListNotificationRequest) GenerateUrl() (string, url.Values, error) {

	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "" , nil, err
	}

	return baseUrl, nil, nil
}

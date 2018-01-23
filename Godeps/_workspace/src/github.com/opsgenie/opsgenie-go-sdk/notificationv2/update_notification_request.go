package notificationv2

import (
	"net/url"
	"errors"
)

// UpdateNotificationRequest is a struct of request to update existing notification rule.
type UpdateNotificationRequest struct {
	*Identifier
	ApiKey           string
	Name             string             `json:"name"`
	Criteria         Criteria           `json:"criteria"`
	NotificationTime []NotificationTime `json:"notificationTime"`
	TimeRestriction  TimeRestriction    `json:"timeRestriction"`
	Schedules        []Schedule         `json:"schedules"`
	Steps            []Step             `json:"steps"`
	Repeat           Repeat             `json:"repeat"`
	Order            int                `json:"order"`
	Enabled          bool               `json:"enabled"`
}

// GetApiKey returns api key.
func (r *UpdateNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *UpdateNotificationRequest) GenerateUrl() (string, url.Values, error) {

	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "" , nil, err
	}

	if r.Identifier.RuleID == "" {
		return "", nil, errors.New("Rule ID should be provided for update action")
	}

	baseUrl += "/" + r.Identifier.RuleID

	return baseUrl, nil, nil
}

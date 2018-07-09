package notificationv2

import (
	"errors"
	"net/url"
)

// CreateNotificationRequest is a struct of request to crate new notification rule.
type CreateNotificationRequest struct {
	*Identifier
	ApiKey           string
	Name             string             `json:"name"`
	ActionType       ActionType         `json:"actionType"`
	Criteria         Criteria           `json:"criteria"`
	NotificationTime []NotificationTime `json:"notificationTime"`
	TimeRestriction  TimeRestriction    `json:"timeRestriction"`
	Schedules        []Schedule         `json:"schedules"`
	Order            int                `json:"order"`
	Steps            []Step             `json:"steps"`
	Repeat           Repeat             `json:"repeat"`
	Enabled          bool               `json:"enabled"`
}

// GetApiKey returns api key.
func (r *CreateNotificationRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *CreateNotificationRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if r.Name == "" {
		return "", nil, errors.New("Name should be provided for create action")
	}

	if r.ActionType == "" {
		return "", nil, errors.New("Action Type should be provided for create action")
	}

	if r.ActionType == ScheduleStartActionType || r.ActionType == ScheduleEndActionType {
		if len(r.NotificationTime) < 1 {
			return "", nil, errors.New("Notification Time should be provided for create action if Action Type selected as Schedule Start or Schedule End")
		}
	}

	return baseUrl, nil, nil
}

package userv2

import "net/url"

// ListUserSchedulesRequest is a request for getting list of user schedules.
type ListUserSchedulesRequest struct {
	*Identifier
	ApiKey string
}

// GetApiKey returns api key.
func (r *ListUserSchedulesRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates API url using specified attributes of identifier.
func (r *ListUserSchedulesRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}
	return baseUrl + "/schedules", params, err;
}

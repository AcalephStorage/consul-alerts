package userv2

import (
	"net/url"
)

// GetUserRequest is a request for getting user.
type GetUserRequest struct {
	*Identifier
	ApiKey string
	ExpandContact bool
}

// GetApiKey returns api key.
func (r *GetUserRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates API url using specified attributes of identifier.
func (request *GetUserRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, params, err := request.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if request.ExpandContact {
		if params == nil {
			params = url.Values{}
		}

		params.Add("expand", "contact")
	}

	return baseUrl, params, nil
}

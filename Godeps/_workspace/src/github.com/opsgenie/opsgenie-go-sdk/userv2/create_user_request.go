package userv2

import (
	"errors"
	"net/url"
)

// CreateUserRequest is a request for creating new user.
type CreateUserRequest struct {
	UserName          string              `json:"username,omitempty"`
	FullName          string              `json:"fullName,omitempty"`
	Role              *UserRole           `json:"role,omitempty"`
	SkypeUsername     string              `json:"skypeUsername,omitempty"`
	UserAddress       UserAddress         `json:"userAddress,omitempty"`
	Tags              []string            `json:"tags,omitempty"`
	Details           map[string][]string `json:"details,omitempty"`
	Timezone          string              `json:"timezone,omitempty"`
	Locale            string              `json:"locale,omitempty"`
	DisableInvitation bool                `json:"invitationDisabled,omitempty"`
	ApiKey            string              `json:"-"`
}

// GenerateUrl generates url to API endpoint.
func (r *CreateUserRequest) GenerateUrl() (string, url.Values, error) {

	if r.UserName == "" {
		return "", nil, errors.New("Username should be provided for create action")
	}

	if r.FullName == "" {
		return "", nil, errors.New("FullName should be provided for create action")
	}

	if r.Role == nil {
		return "", nil, errors.New("Role should be provided for create action")
	}
	return "/v2/users", nil, nil
}

// GetApiKey returns api key.
func (r *CreateUserRequest) GetApiKey() string {
	return r.ApiKey
}

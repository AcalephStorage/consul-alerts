package userv2

import "net/url"

// UpdateUserRequest is a request for updating user.
type UpdateUserRequest struct {
	*Identifier
	Username          string              `json:"username,omitempty"`
	FullName          string              `json:"fullName,omitempty"`
	Role              UserRole            `json:"role,omitempty"`
	SkypeUsername     string              `json:"skypeUsername,omitempty"`
	UserAddress       UserAddress         `json:"userAddress,omitempty"`
	Tags              []string            `json:"tags,omitempty"`
	Details           map[string][]string `json:"details,omitempty"`
	Timezone          string              `json:"timezone,omitempty"`
	Locale            string              `json:"locale,omitempty"`
	DisableInvitation bool                `json:"invitationDisabled,omitempty"`
	ApiKey            string              `json:"-"`
}

// GetApiKey returns api key.
func (r *UpdateUserRequest) GetApiKey() string {
	return r.ApiKey
}

// GenerateUrl generates url to API endpoint.
func (r *UpdateUserRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl, _, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return baseUrl, nil, err
}

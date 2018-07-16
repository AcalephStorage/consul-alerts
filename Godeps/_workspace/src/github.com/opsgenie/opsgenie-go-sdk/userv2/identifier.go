package userv2

import (
	"errors"
	"net/url"
)

// Identifier defined the set of attributes for identification user.
type Identifier struct {
	ID       string `json:"-"`
	Username string `json:"-"`
}

// GenerateUrl generates API url using specified attributes of identifier.
func (request *Identifier) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/users"

	if len(request.ID) > 0 {
		baseUrl += "/" + request.ID
	} else if len(request.Username) > 0 {
		baseUrl += "/" + request.Username
	} else {
		return "", nil, errors.New("username or id of the user should be provided")
	}

	return baseUrl, nil, nil
}

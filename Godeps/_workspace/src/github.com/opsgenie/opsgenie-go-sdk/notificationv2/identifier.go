package notificationv2

import (
	"errors"
	"net/url"
)

// Identifier defined the set of attributes for identification notification.
type Identifier struct {
	UserID         string              `json:"-"`
	Username 	   string 			   `json:"-"`
	RuleID         string              `json:"-"`
}

// GenerateUrl generates API url using specified attributes of identifier.
func (request *Identifier) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/users/"

	if request.UserID != "" {
		 baseUrl += request.UserID + "/notification-rules"
	} else if request.Username != "" {
		baseUrl += request.Username + "/notification-rules"
	} else {
		return  "", nil, errors.New("UserID or Username should be provided")
	}

	return baseUrl, nil, nil
}

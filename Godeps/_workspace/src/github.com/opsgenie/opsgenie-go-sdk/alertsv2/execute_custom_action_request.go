package alertsv2

import (
	"net/url"
	"errors"
)

type ExecuteCustomActionRequest struct {
	*Identifier
	ActionName string `json:"-"`
	User       string `json:"user,omitempty"`
	Source     string `json:"source,omitempty"`
	Note       string `json:"note,omitempty"`
	ApiKey     string `json:"-"`
}

func (r *ExecuteCustomActionRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	if r.ActionName == "" {
		return "", nil, errors.New("ActionName should be provided")
	}
	return path + "/actions/" + r.ActionName, params, err;
}

func (r *ExecuteCustomActionRequest) GetApiKey() string {
	return r.ApiKey
}

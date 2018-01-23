package alertsv2

import "net/url"

type AssignAlertRequest struct {
	*Identifier
	Owner  User `json:"owner,omitempty"`
	User   string `json:"user,omitempty"`
	Source string `json:"source,omitempty"`
	Note   string `json:"note,omitempty"`
	ApiKey string `json:"-"`
}

func (r *AssignAlertRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/assign", params, err;
}

func (r *AssignAlertRequest) GetApiKey() string {
	return r.ApiKey
}

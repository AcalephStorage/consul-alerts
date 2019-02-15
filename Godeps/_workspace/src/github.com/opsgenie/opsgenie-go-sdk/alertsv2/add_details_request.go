package alertsv2

import "net/url"

type AddDetailsRequest struct {
	*Identifier
	Details map[string]string `json:"details,omitempty"`
	User    string            `json:"user,omitempty"`
	Source  string            `json:"source,omitempty"`
	Note    string            `json:"note,omitempty"`
	ApiKey  string            `json:"-"`
}

func (r *AddDetailsRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/details", params, err
}

func (r *AddDetailsRequest) GetApiKey() string {
	return r.ApiKey
}

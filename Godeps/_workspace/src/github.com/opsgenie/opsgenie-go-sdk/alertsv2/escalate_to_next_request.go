package alertsv2

import "net/url"

type EscalateToNextRequest struct {
	*Identifier
	Escalation Escalation `json:"escalation,omitempty"`
	User       string `json:"user,omitempty"`
	Source     string `json:"source,omitempty"`
	Note       string `json:"note,omitempty"`
	ApiKey     string `json:"-"`
}

func (r *EscalateToNextRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/escalate", params, err;
}

func (r *EscalateToNextRequest) GetApiKey() string {
	return r.ApiKey
}

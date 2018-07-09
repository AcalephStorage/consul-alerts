package alertsv2

import url "net/url"

type AddTeamToAlertRequest struct {
	*Identifier
	Team   Team   `json:"team,omitempty"`
	User   string `json:"user,omitempty"`
	Source string `json:"source,omitempty"`
	Note   string `json:"note,omitempty"`
	ApiKey string `json:"-"`
}

func (r *AddTeamToAlertRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/teams", params, err
}

func (r *AddTeamToAlertRequest) GetApiKey() string {
	return r.ApiKey
}

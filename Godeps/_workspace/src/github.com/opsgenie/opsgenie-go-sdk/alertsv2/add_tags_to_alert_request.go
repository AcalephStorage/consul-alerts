package alertsv2

import "net/url"

type AddTagsToAlertRequest struct {
	*Identifier
	Tags   []string `json:"tags,omitempty"`
	User   string `json:"user,omitempty"`
	Source string `json:"source,omitempty"`
	Note   string `json:"note,omitempty"`
	ApiKey string `json:"-"`
}

func (r *AddTagsToAlertRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/tags", params, err;
}

func (r *AddTagsToAlertRequest) GetApiKey() string {
	return r.ApiKey
}

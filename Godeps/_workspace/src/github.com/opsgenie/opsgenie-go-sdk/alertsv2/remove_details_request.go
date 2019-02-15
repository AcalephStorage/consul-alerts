package alertsv2

import (
	"net/url"
	"strings"
)

type RemoveDetailsRequest struct {
	*Identifier
	Keys   []string
	User   string
	Source string
	Note   string
	ApiKey string
}

func (r *RemoveDetailsRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *RemoveDetailsRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if r.User != "" {
		params.Add("user", r.User)
	}

	if r.Source != "" {
		params.Add("source", r.Source)
	}

	if r.Keys != nil {
		params.Add("keys", strings.Join(r.Keys, ","))
	}

	return path + "/details", params, nil
}

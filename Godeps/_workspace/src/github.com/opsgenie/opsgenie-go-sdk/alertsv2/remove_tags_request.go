package alertsv2

import (
	"net/url"
	"strings"
)

type RemoveTagsRequest struct {
	*Identifier
	Tags   []string
	User   string
	Source string
	Note   string
	ApiKey string
}

func (r *RemoveTagsRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *RemoveTagsRequest) GenerateUrl() (string, url.Values, error) {
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

	if r.Tags != nil {
		params.Add("tags", strings.Join(r.Tags, ","))
	}

	return path + "/tags", params, nil
}

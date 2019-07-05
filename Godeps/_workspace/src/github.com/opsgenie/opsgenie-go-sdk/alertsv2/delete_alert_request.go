package alertsv2

import "net/url"

type DeleteAlertRequest struct {
	*Identifier
	User   string
	Source string
	ApiKey string
}

func (r *DeleteAlertRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *DeleteAlertRequest) GenerateUrl() (string, url.Values, error) {
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

	return path, params, err
}

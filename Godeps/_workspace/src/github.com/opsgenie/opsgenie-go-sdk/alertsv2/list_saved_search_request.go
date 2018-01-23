package alertsv2

import url "net/url"

type LisSavedSearchRequest struct {
	ApiKey string
}

func (r *LisSavedSearchRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *LisSavedSearchRequest) GenerateUrl() (string, url.Values, error) {
	return "/v2/alerts/saved-searches", nil, nil
}


package savedsearches

import (
	"errors"
	"net/url"
)

type GetSavedSearchRequest struct {
	ID     string
	Name   string
	ApiKey string
}

func (r *GetSavedSearchRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *GetSavedSearchRequest) GenerateUrl() (string, url.Values, error) {
	path := "/v2/alerts/saved-searches"

	if r.ID != "" {
		return path + "/" + r.ID, nil, nil
	}

	if r.Name != "" {
		params := url.Values{}
		params.Add("identifierType", "name")

		return path + "/" + r.Name, params, nil
	}

	return "", nil, errors.New("ID or Name should be provided")
}

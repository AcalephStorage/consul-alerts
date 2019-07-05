package alertsv2

import (
	"errors"
	"net/url"
)

type Identifier struct {
	ID     string `json:"-"`
	Alias  string `json:"-"`
	TinyID string `json:"-"`
}

func (request *Identifier) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/alerts/"

	if request.ID != "" {
		return baseUrl + request.ID, url.Values{}, nil
	}

	params := url.Values{}

	if request.Alias != "" {
		params.Set("identifierType", "alias")
		return baseUrl + request.Alias, params, nil
	}

	if request.TinyID != "" {
		params.Set("identifierType", "tiny")
		return baseUrl + request.TinyID, params, nil
	}

	return "", nil, errors.New("ID, TinyID or Alias should be provided")
}

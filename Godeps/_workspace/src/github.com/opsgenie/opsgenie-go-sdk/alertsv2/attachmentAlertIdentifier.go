package alertsv2

import (
	"errors"
	"net/url"
)

type AttachmentAlertIdentifier struct {
	ID     string `json:"-"`
	Alias  string `json:"-"`
	TinyID string `json:"-"`
}

func (request *AttachmentAlertIdentifier) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/alerts/"

	if request.ID != "" {
		return baseUrl + request.ID, url.Values{}, nil
	}

	params := url.Values{}

	if request.Alias != "" {
		params.Set("alertIdentifierType", "alias")
		return baseUrl + request.Alias, params, nil
	}

	if request.TinyID != "" {
		params.Set("alertIdentifierType", "tiny")
		return baseUrl + request.TinyID, params, nil
	}

	return "", nil, errors.New("ID, TinyID or Alias should be provided")
}

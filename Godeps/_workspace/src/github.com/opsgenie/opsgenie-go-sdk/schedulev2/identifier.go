package schedulev2

import (
	"net/url"
	"errors"
)

type Identifier struct {
	ID	 	string
	Name 	string
}

func (request *Identifier) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/schedules/"
	params := url.Values{}

	if request.Name != "" {
		baseUrl += url.QueryEscape(request.Name)
		params.Add("identifierType", "name")
	} else if request.ID != "" {
		baseUrl += url.QueryEscape(request.ID)
		params.Add("identifierType", "id")
	} else {
		return "", nil, errors.New("Identifier should be ID or Name")
	}

	return baseUrl, params, nil
}


package alertsv2

import (
	"net/url"
	"strconv"
)

type ListAlertRequest struct {
	Limit                int
	Sort                 SortField
	Offset               int
	Order                Order
	Query                string
	SearchIdentifier     string
	SearchIdentifierType SearchIdentifierType
	ApiKey               string
}

func (r *ListAlertRequest) GetApiKey() string {
	return r.ApiKey;
}

func (request *ListAlertRequest) GenerateUrl() (string, url.Values, error) {
	params := url.Values{}

	if request.Limit != 0 {
		params.Add("limit", strconv.Itoa(request.Limit))
	}

	if request.Sort != "" {
		params.Add("sort", string(request.Sort))
	}

	if request.Offset != 0 {
		params.Add("offset", strconv.Itoa(request.Offset))
	}

	if request.Query != "" {
		params.Add("query", request.Query)
	}

	if request.SearchIdentifier != "" {
		params.Add("searchIdentifier", request.SearchIdentifier)
	}

	if request.SearchIdentifierType != "" {
		params.Add("searchIdentifierType", string(request.SearchIdentifierType))
	}

	return "/v2/alerts", params, nil
}

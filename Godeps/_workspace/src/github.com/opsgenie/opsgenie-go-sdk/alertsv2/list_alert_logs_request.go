package alertsv2

import (
	url "net/url"
	"strconv"
)

type ListAlertLogsRequest struct {
	*Identifier
	Offset    string
	Direction Direction
	Limit     int
	Order     Order
	ApiKey    string
}

func (r *ListAlertLogsRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	if r.Offset != "" {
		params.Add("offset", r.Offset)
	}

	if r.Direction != "" {
		params.Add("direction", string(r.Direction))
	}

	if r.Limit != 0 {
		params.Add("limit", strconv.Itoa(r.Limit))
	}

	if r.Order != "" {
		params.Add("order", string(r.Order))
	}

	return path + "/logs", params, nil
}

func (r *ListAlertLogsRequest) GetApiKey() string {
	return r.ApiKey
}

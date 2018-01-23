package alertsv2

import "net/url"

type ListAlertRecipientsRequest struct {
	*Identifier
	ApiKey string
}

func (r *ListAlertRecipientsRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *ListAlertRecipientsRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return path + "/recipients", params, nil
}

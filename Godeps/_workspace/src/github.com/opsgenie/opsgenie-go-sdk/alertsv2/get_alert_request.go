package alertsv2

type GetAlertRequest struct {
	*Identifier
	ApiKey string
}

func (r *GetAlertRequest) GetApiKey() string {
	return r.ApiKey
}

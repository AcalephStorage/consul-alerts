package alertsv2

type DetailedAlert struct {
	*Alert
	Actions     []string `json:"actions,omitempty"`
	Entity      string `json:"entity,omitempty"`
	Description string `json:"description,omitempty"`
	Details     map[string]string `json:"details,omitempty"`
}

type DetailedAlertResponse struct {
	ResponseMeta
	Alert DetailedAlert `json:"data"`
}

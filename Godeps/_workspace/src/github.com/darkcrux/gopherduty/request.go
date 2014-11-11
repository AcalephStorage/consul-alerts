package gopherduty

import (
	"bytes"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

const endpoint = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"

type pagerDutyRequest struct {
	ServiceKey  string      `json:"service_key"`
	EventType   string      `json:"event_type"`
	IncidentKey string      `json:"incident_key,omitempty"`
	Description string      `json:"description"`
	Client      string      `json:"client,omitempty"`
	ClientUrl   string      `json:"client_url,omitempty"`
	Details     interface{} `json:"details"`
}

func (p *pagerDutyRequest) submit() (pagerResponse *PagerDutyResponse) {
	pagerResponse = &PagerDutyResponse{}

	body, err := json.Marshal(p)
	if err != nil {
		pagerResponse.appendError(err)
		return pagerResponse
	}

	buf := bytes.NewBuffer(body)
	response, err := http.Post(endpoint, "application/json", buf)
	if err != nil {
		pagerResponse.appendError(err)
		return pagerResponse
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		pagerResponse.appendError(err)
		return pagerResponse
	}

	pagerResponse.parse(responseBody)

	return pagerResponse
}

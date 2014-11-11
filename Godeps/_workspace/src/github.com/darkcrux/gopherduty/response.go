package gopherduty

import "encoding/json"

// The response from calling the PagerDuty API. This can contain errors if the API call failed. Also, any errors
// encountered when calling the API is added to the Errors list.
type PagerDutyResponse struct {
	Status      string   `json:"status"`
	Message     string   `json:"message"`
	IncidentKey string   `json:"incident_key,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

// Return the JSON string.
func (p *PagerDutyResponse) String() string {
	resp, _ := json.Marshal(p)
	return string(resp)
}

// Error interface implementation.
func (p *PagerDutyResponse) Error() string {
	return p.String()
}

// Returns true if there are any errors during API call.
func (p *PagerDutyResponse) HasErrors() bool {
	return len(p.Errors) > 0
}

func (p *PagerDutyResponse) parse(rawResponse []byte) {
	if err := json.Unmarshal(rawResponse, p); err != nil {
		p.appendError(err)
	}
}

func (p *PagerDutyResponse) appendError(err error) {
	p.Errors = append(p.Errors, err.Error())
}

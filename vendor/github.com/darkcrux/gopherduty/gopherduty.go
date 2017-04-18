// A simple Go client for PagerDuty's API. This includes the trigger, acknowledge, and
// resolve event types. This also includes a retry feature when sending to PagerDuty
// fails.
package gopherduty

import (
	"log"
	"math"
	"time"
)

const (
	eventTrigger     = "trigger"
	eventAcknowledge = "acknowledge"
	eventResolve     = "resolve"
)

func init() {
	log.SetPrefix("[ PagerDuty Client ] ")
}

// PagerDuty requires a Service Key to work. API call can be retried if MaxRetry is set to > 1. This retries the
// request with an exponential delay for each retry.
type PagerDuty struct {
	ServiceKey        string // The Service key needed to access PagerDuty.
	MaxRetry          int    // Maximum API call retries. Defaults to 0.
	RetryBaseInterval int    // Starting delay for a retry in seconds. Defaults to 10.
	retries           int
}

// Convenience method to create a new PagerDuty struct.
func NewClient(serviceKey string) *PagerDuty {
	return &PagerDuty{
		ServiceKey: serviceKey,
	}
}

// Send a TRIGGER event. The incidentKey may be left empty and PagerDuty will generate one.
func (p *PagerDuty) Trigger(incidentKey, description, client, clientUrl string, details interface{}) *PagerDutyResponse {
	log.Println("Sending TRIGGER event")
	return p.doRequest(eventTrigger, incidentKey, description, client, clientUrl, details)
}

// Send an ACKNOWLEDGE event.
func (p *PagerDuty) Acknowledge(incidentKey, description string, details interface{}) *PagerDutyResponse {
	log.Println("Sending ACKENOWLEDGE event")
	return p.doRequest(eventAcknowledge, incidentKey, description, "", "", details)
}

// Send a RESOLVE event.
func (p *PagerDuty) Resolve(incidentKey, description string, details interface{}) *PagerDutyResponse {
	log.Println("Sending RESOLVE event")
	return p.doRequest(eventResolve, incidentKey, description, "", "", details)
}

func (p *PagerDuty) doRequest(eventType, incidentKey, description, client, clientUrl string, details interface{}) *PagerDutyResponse {
	request := &pagerDutyRequest{
		ServiceKey:  p.ServiceKey,
		EventType:   eventType,
		IncidentKey: incidentKey,
		Description: description,
		Client:      client,
		ClientUrl:   clientUrl,
		Details:     details,
	}

	response := request.submit()
	if response.HasErrors() && p.retries < p.MaxRetry {
		p.delayRetry()
		p.retries++
		response = p.doRequest(eventType, incidentKey, description, client, clientUrl, details)
	}
	p.retries = 0
	return response
}

func (p *PagerDuty) delayRetry() {
	interval := float64(p.RetryBaseInterval)
	if interval == 0 {
		interval = 10
	}
	delay := math.Pow(2, float64(p.retries)) * interval
	duration := time.Duration(delay) * time.Second

	log.Printf("Retrying in %v...\n", duration)
	time.Sleep(duration)
}

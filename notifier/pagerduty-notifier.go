package notifier

import (
	log "github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/darkcrux/gopherduty"
)

const defaultRetryBaseInterval = 30

type PagerDutyNotifier struct {
	Enabled           bool
	ServiceKey        string `json:"service-key"`
	ClientName        string `json:"client-name"`
	ClientUrl         string `json:"client-url"`
	MaxRetry          int    `json:"max-retry"`
	RetryBaseInterval int    `json:"retry-base-interval"'`
}

// NotifierName provides name for notifier selection
func (pd *PagerDutyNotifier) NotifierName() string {
	return "pagerduty"
}

func (pd *PagerDutyNotifier) Copy() Notifier {
	notifier := *pd
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (pd *PagerDutyNotifier) Notify(messages Messages) bool {

	client := gopherduty.NewClient(pd.ServiceKey)

	if pd.MaxRetry != 0 {
		client.MaxRetry = pd.MaxRetry

		if pd.RetryBaseInterval != 0 {
			client.RetryBaseInterval = pd.RetryBaseInterval
		} else {
			client.RetryBaseInterval = defaultRetryBaseInterval
		}
	}

	result := true

	for _, message := range messages {
		incidentKey := message.Node
		if message.ServiceId != "" {
			incidentKey += ":" + message.ServiceId
		}
		incidentKey += ":" + message.CheckId
		subject := message.Node
		if message.Service != "" {
			subject += ":" + message.Service
		}
		if message.Check != "" {
			subject += ":" + message.Check
		}
		var response *gopherduty.PagerDutyResponse
		switch {
		case message.IsPassing():
			description := subject + " is now HEALTHY"
			response = client.Resolve(incidentKey, description, message)
		case message.IsWarning():
			description := subject + " is UNSTABLE"
			response = client.Trigger(incidentKey, description, pd.ClientName, pd.ClientUrl, message)
		case message.IsCritical():
			description := subject + " is CRITICAL"
			response = client.Trigger(incidentKey, description, pd.ClientName, pd.ClientUrl, message)
		}

		if response.HasErrors() {
			for _, err := range response.Errors {
				log.Printf("Error sending %s notification to pagerduty: %s\n", incidentKey, err)
			}
			result = false
		}
	}

	log.Println("PagerDuty notification complete")
	return result
}

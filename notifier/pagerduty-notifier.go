package notifier

import (
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/darkcrux/gopherduty"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type PagerDutyNotifier struct {
	Enabled    bool
	ServiceKey string `json:"service-key"`
	ClientName string `json:"client-name"`
	ClientUrl  string `json:"client-url"`
}

// NotifierName provides name for notifier selection
func (pd *PagerDutyNotifier) NotifierName() string {
	return "pagerduty"
}

//Notify sends messages to the endpoint notifier
func (pd *PagerDutyNotifier) Notify(messages Messages) bool {

	client := gopherduty.NewClient(pd.ServiceKey)

	result := true

	for _, message := range messages {
		incidentKey := message.Node
		if message.ServiceId != "" {
			incidentKey += ":" + message.ServiceId
		}
		incidentKey += ":" + message.CheckId
		var response *gopherduty.PagerDutyResponse
		switch {
		case message.IsPassing():
			description := incidentKey + " is now HEALTHY"
			response = client.Resolve(incidentKey, description, message)
		case message.IsWarning():
			description := incidentKey + " is UNSTABLE"
			response = client.Trigger(incidentKey, description, pd.ClientName, pd.ClientUrl, message)
		case message.IsCritical():
			description := incidentKey + " is CRITICAL"
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

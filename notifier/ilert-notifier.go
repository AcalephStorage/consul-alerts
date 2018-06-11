package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const apiEndpoint string = "https://ilertnow.com/api/v1/events"

type ILertNotifier struct {
	Enabled                bool
	ApiKey                 string `json:"api-key"`
	IncidentKeyIncludeHost bool   `json:"incident-key-include-host"`
}

type iLertEvent struct {
	ApiKey      string `json:"apiKey"`
	EventType   string `json:"eventType"`
	Summary     string `json:"summary"`
	Details     string `json:"details"`
	IncidentKey string `json:"incidentKey"`
}

// NotifierName provides name for notifier selection
func (il *ILertNotifier) NotifierName() string {
	return "ilert"
}

func (il *ILertNotifier) Copy() Notifier {
	notifier := *il
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (il *ILertNotifier) Notify(messages Messages) bool {
	result := true

	for _, message := range messages {
		var eventType string
		var summary string

		ik := il.incidentKey(message)

		switch {
		case message.IsPassing():
			summary = ik + " is now HEALTHY"
			eventType = "RESOLVE"
		case message.IsWarning():
			// iLert does not support warning state
			continue
		case message.IsCritical():
			summary = ik + " is CRITICAL"
			eventType = "ALERT"
		}

		if err := il.sendEvent(eventType, summary, message.Output, ik); err != nil {
			log.Error("Problem while sending iLert event:", err)
			result = false
		}
	}

	log.Println("iLert notification complete")
	return result
}

//sendEvent builds the event JSON and sends it to the iLert API
func (il *ILertNotifier) sendEvent(eventType, summary, details, incidentKey string) error {
	event := iLertEvent{
		ApiKey:      il.ApiKey,
		EventType:   eventType,
		Summary:     summary,
		Details:     details,
		IncidentKey: incidentKey,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	log.Debugf("struct = %+v, json = %s", event, string(body))

	res, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("unexpected HTTP status code: %d (%s)", res.StatusCode, string(body)))
	}

	return nil
}

func (il *ILertNotifier) incidentKey(message Message) string {
	if il.IncidentKeyIncludeHost {
		return fmt.Sprintf("%s:%s:%s", message.Node, message.Service, message.Check)
	} else {
		return fmt.Sprintf("%s:%s", message.Service, message.Check)
	}
}

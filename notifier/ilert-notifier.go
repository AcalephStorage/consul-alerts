package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const apiEndpoint string = "https://ilertnow.com/api/v1/events"

type ILertNotifier struct {
	ApiKey              string `json:"api-key"`
	Enabled             bool
	IncidentKeyTemplate string `json:"incident-key-template"`

	incidentKeyTemplateCompiled *template.Template
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

func (il *ILertNotifier) toILertEvents(messages Messages) []iLertEvent {
	iLertEvents := make([]iLertEvent, 0)

	for _, message := range messages {
		var eventType string
		var summary string

		ik, err := il.incidentKey(message)
		if err != nil {
			log.Error("Failed to create an iLert event: ", err)
			continue
		}

		switch {
		case message.IsPassing():
			summary = ik + " is now HEALTHY"
			eventType = "RESOLVE"
		case message.IsWarning():
			summary = ik + " is WARNING"
			eventType = "RESOLVE"
		case message.IsCritical():
			summary = ik + " is CRITICAL"
			eventType = "ALERT"
		}

		iLertEvents = append(iLertEvents, iLertEvent{
			ApiKey:      il.ApiKey,
			EventType:   eventType,
			Summary:     summary,
			Details:     message.Output,
			IncidentKey: ik,
		})
	}

	return iLertEvents
}

//Notify sends messages to the endpoint notifier
func (il *ILertNotifier) Notify(messages Messages) bool {
	result := true

	for _, iLertEvent := range il.toILertEvents(messages) {
		if err := il.sendEvent(iLertEvent); err != nil {
			log.Error("Problem while sending iLert event: ", err)
			result = false
		}
	}

	log.Println("iLert notification complete")
	return result
}

//sendEvent builds the event JSON and sends it to the iLert API
func (il *ILertNotifier) sendEvent(event iLertEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	res, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New(fmt.Sprintf("Unexpected HTTP status code: %d (%s)", res.StatusCode, string(body)))
	}

	return nil
}

func (il *ILertNotifier) incidentKey(message Message) (string, error) {
	if il.incidentKeyTemplateCompiled == nil {
		il.incidentKeyTemplateCompiled = template.Must(template.New("IncidentKey").Parse(il.IncidentKeyTemplate))
	}

	var buff bytes.Buffer

	err := il.incidentKeyTemplateCompiled.ExecuteTemplate(&buff, "IncidentKey", message)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to render incident key: %s", err))
	}

	return buff.String(), nil
}

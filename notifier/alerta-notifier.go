package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"text/template"

	log "github.com/Difrex/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type AlertaNotifier struct {
	ClusterName string           `json:"-"`
	Domain      string           `json:"-"`
	Url         string           `json:"-"`
	Token       string           `json:"-"`
	Text        string           `json:"text,omitempty"`
	Environment string           `json:"environment"`
	Resource    string           `json:"resource"`
	Event       string           `json:"event"`
	Enabled     bool             `json:"-"`
	Type        string           `json:"type"`
	Origin      string           `json:"origin"`
	Status      string           `json:"status"`
	Service     []string         `json:"service"`
	Attributes  AlertaAttributes `json:"attributes"`
	Severity    string           `json:"severity"`
	Events      []AlertaNotifier `json:"-"`
}

type AlertaAttributes struct {
	Link string `json:"link"`
	Ack  string `json:"ack"`
}

type TmplMsg struct {
	Notifier *AlertaNotifier
	Msg      Message
}

func tpl(t string, msg TmplMsg) (string, error) {
	tmpl, err := template.New("template").Parse(t)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	var s bytes.Buffer
	err = tmpl.Execute(&s, msg)
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

// populateDefaults set default values
func (n *AlertaNotifier) populate(message Messages) {
	for _, msg := range message {

		log.Debug("Alerta notifier: received message: ", msg)

		myEvent := *n
		if myEvent.Type == "" {
			myEvent.Type = "consul-alerts"
		}

		if myEvent.Origin == "" {
			myEvent.Origin = strings.Join([]string{myEvent.ClusterName, myEvent.Domain}, ".")
		}

		if myEvent.Environment == "" {
			myEvent.Environment = "Production"
		}

		msg.Check = strings.ToLower(strings.Replace(msg.Check, " ", "_", 1))
		t := TmplMsg{
			Notifier: n,
			Msg:      msg,
		}

		event, err := tpl(myEvent.Event, t)
		if err != nil {
			log.Error(err.Error())
			return
		}
		myEvent.Event = event

		if myEvent.Resource == "" {
			myEvent.Resource = msg.Node
		} else {
			myEvent.Resource, err = tpl(myEvent.Resource, t)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}

		// Populate services
		myEvent.Service = append(myEvent.Service, msg.Service)
		if myEvent.Attributes.Link != "" {
			link, err := tpl(myEvent.Attributes.Link, t)
			if err != nil {
				log.Error(err.Error())
				return
			}
			myEvent.Attributes.Link = link
		}

		if msg.IsCritical() {
			myEvent.Severity = "major"
			myEvent.Status = "open"
		}

		if msg.IsWarning() {
			myEvent.Severity = "warning"
			myEvent.Status = "open"
		}

		if msg.IsPassing() {
			myEvent.Severity = "ok"
			myEvent.Status = "closed"
		}
		n.Events = append(n.Events, myEvent)
	}
}

// NotifierName provides name for notifier selection
func (n *AlertaNotifier) NotifierName() string {
	return "alerta"
}

func (n *AlertaNotifier) Copy() Notifier {
	notifier := *n
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (n *AlertaNotifier) Notify(messages Messages) bool {
	return n.notifySimple(messages)
}

func (n *AlertaNotifier) notifySimple(messages Messages) bool {
	overallStatus, pass, warn, fail := messages.Summary()
	text := fmt.Sprintf(header, n.ClusterName, overallStatus, fail, warn, pass)
	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
	}
	n.Text = text
	n.populate(messages)

	return n.postToAlerta()
}

func (n *AlertaNotifier) postToAlerta() bool {
	for _, event := range n.Events {
		jsonData, err := json.Marshal(event)
		if err != nil {
			log.Println("Unable to marshal Alerta payload:", err)
			return false
		}

		log.Debugf("struct = %+v, payload = %s", n, string(jsonData))

		client := http.Client{}
		b := bytes.NewBufferString(string(jsonData))
		req, err := http.NewRequest("POST", n.Url, b)
		if err != nil {
			log.Println("Unable to send data to Alerta: ", err.Error())
			return false
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", n.Token)

		res, err := client.Do(req)
		if err != nil {
			log.Println("Unable to send data to Alerta:", err)
			return false
		}
		defer res.Body.Close()

		statusCode := res.StatusCode
		if statusCode > 202 {
			body, _ := ioutil.ReadAll(res.Body)
			log.Println("Unable to notify Alerta: ", string(body))
			return false
		}
		log.Println("Alerta notification sent.")
	}
	return true
}

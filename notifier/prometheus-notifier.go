package notifier

import (
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"html/template"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type PrometheusNotifier struct {
	Enabled     bool
	ClusterName string            `json:"cluster-name"`
	BaseURLs    []string          `json:"base-urls"`
	Endpoint    string            `json:"endpoint"`
	Payload     map[string]string `json:"payload"`
}

type TemplatePayloadData struct {
	Node      string
	Service   string
	Check     string
	Status    string
	Output    string
	Notes     string
	Timestamp string
}

// NotifierName provides name for notifier selection
func (notifier *PrometheusNotifier) NotifierName() string {
	return "prometheus"
}

func (notifier *PrometheusNotifier) Copy() Notifier {
	n := *notifier
	return &n
}

func renderPayload(t TemplatePayloadData, templateFile string, defaultTemplate string) (string, error) {
	var tmpl *template.Template
	var err error
	if templateFile == "" {
		tmpl, err = template.New("base").Parse(defaultTemplate)
	} else {
		tmpl, err = template.ParseFiles(templateFile)
	}

	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, t); err != nil {
		return "", err
	}

	return body.String(), nil
}

//Notify sends messages to the endpoint notifier
func (notifier *PrometheusNotifier) Notify(messages Messages) bool {
	var values []map[string]map[string]string

	for _, m := range messages {
		value := map[string]string{}
		t := TemplatePayloadData{
			Node:      m.Node,
			Service:   m.Service,
			Check:     m.Check,
			Status:    m.Status,
			Output:    m.Output,
			Notes:     m.Notes,
			Timestamp: m.Timestamp.Format("2006-01-02T15:04:05-0700"),
		}

		for payloadKey, payloadVal := range notifier.Payload {
			data, err := renderPayload(t, "", payloadVal)
			if err != nil {
				log.Println("Error rendering template: ", err)
				return false
			}
			value[payloadKey] = string(data)
		}

		values = append(values, map[string]map[string]string{"labels": value})
	}

	requestBody, err := json.Marshal(values)
	if err != nil {
		log.Println("Unable to encode POST data")
		return false
	}

	c := make(chan bool)
	defer close(c)
	for _, bu := range notifier.BaseURLs {
		endpoint := fmt.Sprintf("%s%s", bu, notifier.Endpoint)

		// Channel senders. Logging the result where needed, and sending status back
		go func() {
			if res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody)); err != nil {
				log.Printf(fmt.Sprintf("Unable to send data to Prometheus server (%s): %s", endpoint, err))
				c <- false
			} else {
				defer res.Body.Close()
				statusCode := res.StatusCode

				if statusCode != 200 {
					body, _ := ioutil.ReadAll(res.Body)
					log.Printf(fmt.Sprintf("Unable to notify Prometheus server (%s): %s", endpoint, string(body)))
					c <- false
				} else {
					log.Printf(fmt.Sprintf("Notification sent to Prometheus server (%s).", endpoint))
					c <- true
				}
			}
		}()
	}

	// Channel receiver. Making sure to return the final result in bool
	for i := 0; i < len(notifier.BaseURLs); i++ {
		select {
		case r := <- c:
			if (! r) {
				return false
			}
		}
	}

	return true
}

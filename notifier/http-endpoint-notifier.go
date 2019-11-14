package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type HttpEndpointNotifier struct {
	Enabled     bool
	ClusterName string            `json:"cluster-name"`
	BaseURL     string            `json:"base-url"`
	Endpoint    string            `json:"endpoint"`
	Payload     map[string]string `json:"payload"`
	Token       string            `json:"token"`
}

// NotifierName provides name for notifier selection
func (notifier *HttpEndpointNotifier) NotifierName() string {
	return "http-endpoint"
}

func (notifier *HttpEndpointNotifier) Copy() Notifier {
	n := *notifier
	return &n
}

//Notify sends messages to the endpoint notifier
func (notifier *HttpEndpointNotifier) Notify(messages Messages) bool {
	overallStatus, pass, warn, fail := messages.Summary()
	t := TemplateData{
		ClusterName:  notifier.ClusterName,
		SystemStatus: overallStatus,
		FailCount:    fail,
		WarnCount:    warn,
		PassCount:    pass,
		Nodes:        mapByNodes(messages),
	}
	values := map[string]string{}

	for key, val := range notifier.Payload {
		data, err := renderTemplate(t, "", val)
		if err != nil {
			log.Println("Error rendering template: ", err)
			return false
		}
		values[key] = string(data)
	}

	requestBody, err := json.Marshal(values)
	if err != nil {
		log.Println("Unable to encode POST data")
		return false
	}

	endpoint := fmt.Sprintf("%s%s", notifier.BaseURL, notifier.Endpoint)
	if res, err := sendHTTPWithHeaderAuthorization(endpoint, "application/json", notifier.Token, bytes.NewBuffer(requestBody)); err != nil {
		log.Println("Unable to send data to HTTP endpoint:", err)
		return false
	} else {
		defer res.Body.Close()
		statusCode := res.StatusCode
		if statusCode != 200 {
			body, _ := ioutil.ReadAll(res.Body)
			log.Println("Unable to notify HTTP endpoint:", string(body))
			return false
		} else {
			log.Println("Notification sent to HTTP endpoint.")
			return true
		}
	}

}

// sendHTTPWithHeaderAuthorization ...
func sendHTTPWithHeaderAuthorization(url, contentType string, token string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", contentType)
	if len(token) > 0 {
		req.Header.Set("Authorization", token)

	}
	return client.Do(req)
}

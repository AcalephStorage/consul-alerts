package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type MattermostWebhookNotifier struct {
	ClusterName string `json:"cluster_name"`
	Url         string `json:"-"`
	Channel     string `json:"channel"`
	Username    string `json:"username"`
	IconUrl     string `json:"icon_url"`
	Text        string `json:"text,omitempty"`
	Detailed    bool   `json:"-"`
	Enabled     bool   `json:"enabled"`
}

// NotifierName provides name for notifier selection
func (n *MattermostWebhookNotifier) NotifierName() string {
	return "mattermost-webhook"
}

func (n *MattermostWebhookNotifier) Copy() Notifier {
	notifier := *n
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (n *MattermostWebhookNotifier) Notify(messages Messages) bool {
	return n.notifySimple(messages)
}

func (n *MattermostWebhookNotifier) notifySimple(messages Messages) bool {
	overallStatus, pass, warn, fail := messages.Summary()
	text := fmt.Sprintf(header, n.ClusterName, overallStatus, fail, warn, pass)
	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
	}
	n.Text = text
	return n.postToMattermostWebhook()
}

func (n *MattermostWebhookNotifier) postToMattermostWebhook() bool {
	jsonData, err := json.Marshal(n)
	if err != nil {
		log.Println("Unable to marshal Mattermost payload:", err)
		return false
	}

	values := url.Values{}
	values.Set("payload", string(jsonData))
	data := values.Encode()
	log.Debugf("struct = %+v, payload = %s", n, string(data))

	b := bytes.NewBufferString(data)
	res, err := http.Post(n.Url, "application/x-www-form-urlencoded", b)
	if err != nil {
		log.Println("Unable to send data to Mattermost:", err)
		return false
	}
	defer res.Body.Close()
	statusCode := res.StatusCode
	if statusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		log.Println("Unable to notify Mattermost:", string(body))
		return false
	}
	log.Println("Mattermost notification sent.")
	return true
}

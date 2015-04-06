package notifier

import (
	"bytes"
	"fmt"

	"io/ioutil"

	"encoding/json"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const header = `%s is %s.

Fail: %d, Warn: %d, Pass: %d
`

type SlackNotifier struct {
	ClusterName string `json:"-"`
	Url         string `json:"-"`
	Channel     string `json:"channel"`
	Username    string `json:"username"`
	IconUrl     string `json:"icon_url"`
	IconEmoji   string `json:"icon_emoji"`
	Text        string `json:"text"`
	Detailed    bool   `json:"-"`
}

func (slack *SlackNotifier) Notify(messages Messages) bool {
	if slack.Detailed {
		return slack.notifyDetailed(messages)
	} else {
		return slack.notifySimple(messages)
	}
}

func (slack *SlackNotifier) notifySimple(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf(header, slack.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
	}

	slack.Text = text

	data, err := json.Marshal(slack)
	if err != nil {
		log.Println("Unable to marshal slack payload:", err)
		return false
	}

	b := bytes.NewBuffer(data)
	if res, err := http.Post(slack.Url, "application/json", b); err != nil {
		log.Println("Unable to send data to slack:", err)
		return false
	} else {
		defer res.Body.Close()
		statusCode := res.StatusCode
		if statusCode != 200 {
			body, _ := ioutil.ReadAll(res.Body)
			log.Println("Unable to notify slack:", string(body))
			return false
		} else {
			log.Println("Slack notification sent.")
			return true
		}
	}

}

func (slack *SlackNotifier) notifyDetailed(messages Messages) bool {
	// TBD
}

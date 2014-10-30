package notifier

import (
	"bytes"
	"fmt"

	"encoding/json"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

const SlackUrl = "https://%s.slack.com/services/hooks/incoming-webhook?token=%s"
const header = `%s is %s.

Fail: %d, Warn: %d, Pass: %d
`

type SlackNotifier struct {
	ClusterName string `json:"-"`
	Team        string `json:"-"`
	Token       string `json:"-"`
	Channel     string `json:"channel"`
	Username    string `json:"username"`
	IconUrl     string `json:"icon_url"`
	IconEmoji   string `json:"icon_emoji"`
	Text        string `json:"text"`
}

func (slack *SlackNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf(header, slack.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
	}

	url := fmt.Sprintf(SlackUrl, slack.Team, slack.Token)
	slack.Text = text

	data, err := json.Marshal(slack)
	if err != nil {
		log.Println("Unable to marshal slack payload:", err)
		return false
	}

	b := bytes.NewBuffer(data)
	if res, err := http.Post(url, "application/json", b); err != nil {
		log.Println("Unable to send data to slack:", err)
		return false
	} else {
		res.Body.Close()
		log.Println("Slack notification sent.")
		return true
	}

}

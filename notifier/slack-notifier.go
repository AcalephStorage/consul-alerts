package notifier

import (
	"bytes"
	"fmt"
	"strings"

	"io/ioutil"

	"encoding/json"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type SlackNotifier struct {
	ClusterName string       `json:"-"`
	Url         string       `json:"-"`
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	IconUrl     string       `json:"icon_url"`
	IconEmoji   string       `json:"icon_emoji"`
	Text        string       `json:"text,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
	Detailed    bool         `json:"-"`
	Enabled     bool
}

type attachment struct {
	Color    string   `json:"color"`
	Title    string   `json:"title"`
	Pretext  string   `json:"pretext"`
	Text     string   `json:"text"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

// NotifierName provides name for notifier selection
func (slack *SlackNotifier) NotifierName() string {
	return "slack"
}

//Notify sends messages to the endpoint notifier
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

	return slack.postToSlack()

}

func (slack *SlackNotifier) notifyDetailed(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	var emoji, color string
	switch overallStatus {
	case SYSTEM_HEALTHY:
		emoji = ":white_check_mark:"
		color = "good"
	case SYSTEM_UNSTABLE:
		emoji = ":question:"
		color = "warning"
	case SYSTEM_CRITICAL:
		emoji = ":x:"
		color = "danger"
	default:
		emoji = ":question:"
	}
	title := "Consul monitoring report"
	pretext := fmt.Sprintf("%s %s is *%s*", emoji, slack.ClusterName, overallStatus)

	detailedBody := ""
	detailedBody += fmt.Sprintf("*Changes:* Fail = %d, Warn = %d, Pass = %d", fail, warn, pass)
	detailedBody += fmt.Sprintf("\n")

	for _, message := range messages {
		detailedBody += fmt.Sprintf("\n*[%s:%s]* %s is *%s.*", message.Node, message.Service, message.Check, message.Status)
		var msg = strings.TrimSpace(message.Output)
		if len(msg) != 0 {
			detailedBody += fmt.Sprintf("\n`%s`", msg)
		}
	}

	a := attachment{
		Color:    color,
		Title:    title,
		Pretext:  pretext,
		Text:     detailedBody,
		MrkdwnIn: []string{"text", "pretext"},
	}
	slack.Attachments = []attachment{a}

	return slack.postToSlack()

}

func (slack *SlackNotifier) postToSlack() bool {

	data, err := json.Marshal(slack)
	if err != nil {
		log.Println("Unable to marshal slack payload:", err)
		return false
	}
	log.Debugf("struct = %+v, json = %s", slack, string(data))

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

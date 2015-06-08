package notifier

import (
	"fmt"
	"net/url"

	"github.com/tbruyelle/hipchat-go/hipchat"

	log "github.com/Sirupsen/logrus"
)

type HipChatNotifier struct {
	ClusterName string
	RoomId      string
	AuthToken   string
	BaseURL     string
}

func (notifier *HipChatNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf(header, notifier.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
	}

	level := "green"
	if fail > 0 {
		level = "red"
	} else if warn > 0 {
		level = "yellow"
	}

	client := hipchat.NewClient(notifier.AuthToken)
	if notifier.BaseURL != "" {
		url, err := url.Parse(notifier.BaseURL)
		if err != nil {
			log.Printf("Error parsing hipchat base url: %s\n", err)
		}
		client.BaseURL = url
	}

	notifRq := &hipchat.NotificationRequest{
		Message: text,
		Color:   level,
		Notify:  true,
	}
	resp, err := client.Room.Notification(notifier.RoomId, notifRq)
	if err != nil {
		log.Printf("Error sending notification to hipchat: %s\n", err)
		log.Printf("Server returns %+v\n", resp)
		return false
	}

	return true
}

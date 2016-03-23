package notifier

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/tbruyelle/hipchat-go/hipchat"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type HipChatNotifier struct {
	ClusterName string
	RoomId      string
	AuthToken   string
	BaseURL     string
	From        string
	NotifName   string
}

// NotifierName provides name for notifier selection
func (notifier *HipChatNotifier) NotifierName() string {
	return notifier.NotifName
}

//Notify sends messages to the endpoint notifier
func (notifier *HipChatNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf("%s is <STRONG>%s</STRONG>. Fail: %d, Warn: %d, Pass: %d",
                        notifier.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("<BR><CODE>%s</CODE>:%s:%s is <STRONG>%s</STRONG>.",
                            message.Node, html.EscapeString(message.Service), html.EscapeString(message.Check), message.Status)
		text += fmt.Sprintf("<BR>%s", strings.Replace(html.EscapeString(message.Output), "\n", "<BR>", -1);
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

	from := ""
	if notifier.From != "" {
		from = notifier.From
	}

	notifRq := &hipchat.NotificationRequest{
		From:    from,
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

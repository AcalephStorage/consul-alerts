package notifier

import (
    "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/andybons/hipchat"

    log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type HipChatNotifier struct {
    ServiceKey  string
    Room        string
    Sender      string
    Color       string
}

func (hc *HipChatNotifier) Notify(messages Messages) bool {

    client := hipchat.Client{AuthToken: hc.ServiceKey}

    result := true

    for _, message := range messages {
        req := hipchat.MessageRequest{
        RoomId: hc.Room,
        From: hc.Sender,
        switch {
        case message.IsPassing():
            description := message.ServiceId + " is now HEALTHY"
        case message.IsWarning():
            description := message.ServiceId + " is UNSTABLE"
        case message.IsCritical():
            description := incidentKey + " is CRITICAL"
        }
        Message: description,        
        Color: hc.Color,
        MessageFormat: hipchat.FormatText,
        Notify: true,
    }

    if err := c.PostMessage(req); err != nil {
        log.Printf("Error sending %s notification to hipchat: %s\n", message.ServiceId, err)
    }

    log.Println("HipChat notification complete")
    return result
}

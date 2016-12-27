// Package notifier manages notifications for consul-alerts
package notifier

import (
	"time"
)

const (
	SYSTEM_HEALTHY  string = "HEALTHY"
	SYSTEM_UNSTABLE string = "UNSTABLE"
	SYSTEM_CRITICAL string = "CRITICAL"
)

const header = `%s is %s.

Fail: %d, Warn: %d, Pass: %d
`

type Message struct {
	Node         string
	ServiceId    string
	Service      string
	CheckId      string
	Check        string
	Status       string
	Output       string
	Notes        string
	Interval     int
	RmdCheck     time.Time
	NotifList    map[string]bool
	VarOverrides Notifiers
	Timestamp    time.Time
}

type Messages []Message

type Notifier interface {
	Notify(alerts Messages) bool
	NotifierName() string
}

type Notifiers struct {
	Email     *EmailNotifier     `json:"email"`
	Log       *LogNotifier       `json:"log"`
	Influxdb  *InfluxdbNotifier  `json:"influxdb"`
	Slack     *SlackNotifier     `json:"slack"`
	PagerDuty *PagerDutyNotifier `json:"pagerduty"`
	HipChat   *HipChatNotifier   `json:"hipchat"`
	OpsGenie  *OpsGenieNotifier  `json:"opsgenie"`
	AwsSns    *AwsSnsNotifier    `json:"awssns"`
	VictorOps *VictorOpsNotifier `json:"victorops"`
	Custom    []string           `json:"custom"`
}

func (n Notifiers) GetNotifier(name string) (Notifier, bool) {
	switch name {
	case n.Email.NotifierName():
		return n.Email, true
	case n.Log.NotifierName():
		return n.Log, true
	case n.Influxdb.NotifierName():
		return n.Influxdb, true
	case n.Slack.NotifierName():
		return n.Slack, true
	case n.HipChat.NotifierName():
		return n.HipChat, true
	case n.PagerDuty.NotifierName():
		return n.PagerDuty, true
	case n.OpsGenie.NotifierName():
		return n.OpsGenie, true
	case n.AwsSns.NotifierName():
		return n.AwsSns, true
	case n.VictorOps.NotifierName():
		return n.VictorOps, true
	default:
		return nil, false
	}
}

func (m Message) IsCritical() bool {
	return m.Status == "critical"
}

func (m Message) IsWarning() bool {
	return m.Status == "warning"
}

func (m Message) IsPassing() bool {
	return m.Status == "passing"
}

func (m Messages) Summary() (overallStatus string, pass, warn, fail int) {
	hasCritical := false
	hasWarnings := false
	for _, message := range m {
		switch {
		case message.IsCritical():
			hasCritical = true
			fail++
		case message.IsWarning():
			hasWarnings = true
			warn++
		case message.IsPassing():
			pass++
		}
	}
	if hasCritical {
		overallStatus = SYSTEM_CRITICAL
	} else if hasWarnings {
		overallStatus = SYSTEM_UNSTABLE
	} else {
		overallStatus = SYSTEM_HEALTHY
	}
	return
}

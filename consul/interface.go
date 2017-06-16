package consul

import (
	"time"

	notifier "github.com/AcalephStorage/consul-alerts/notifier"
)

// Event data from consul
type Event struct {
	ID            string
	Name          string
	Payload       []byte
	NodeFilter    string
	ServiceFilter string
	TagFilter     string
	Version       uint
	LTime         uint
}

type Check struct {
	Node        string
	CheckID     string
	Name        string
	Status      string
	Notes       string
	Output      string
	ServiceID   string
	ServiceName string
}

type ConsulAlertConfig struct {
	Checks    *ChecksConfig
	Events    *EventsConfig
	Notifiers *notifier.Notifiers
}

type ChecksConfig struct {
	Enabled         bool
	ChangeThreshold int
}

type EventsConfig struct {
	Enabled  bool
	Handlers []string
}

type Status struct {
	Current          string
	CurrentTimestamp time.Time
	Pending          string
	PendingTimestamp time.Time
	HealthCheck      *Check
	ForNotification  bool
}

// ProfileInfo is for reading in JSON from profile keys
type ProfileInfo struct {
	Interval     int
	NotifList    map[string]bool
	VarOverrides notifier.Notifiers
}

// Consul interface provides access to consul client
type Consul interface {
	LoadConfig()

	EventsEnabled() bool
	ChecksEnabled() bool
	EventHandlers(eventName string) []string

	EmailNotifier() *notifier.EmailNotifier
	LogNotifier() *notifier.LogNotifier
	InfluxdbNotifier() *notifier.InfluxdbNotifier
	SlackNotifier() *notifier.SlackNotifier
	MattermostNotifier() *notifier.MattermostNotifier
	PagerDutyNotifier() *notifier.PagerDutyNotifier
	HipChatNotifier() *notifier.HipChatNotifier
	OpsGenieNotifier() *notifier.OpsGenieNotifier
	TelegramNotifier() *notifier.TelegramNotifier
	AwsSnsNotifier() *notifier.AwsSnsNotifier
	VictorOpsNotifier() *notifier.VictorOpsNotifier

	CheckChangeThreshold() int
	UpdateCheckData()
	NewAlerts() []Check
	NewAlertsWithFilter(node string, service string, checkId string, statuses []string, ignoreBlacklist bool) []Check

	IsBlacklisted(check *Check) bool

	CustomNotifiers() map[string]string

	CheckStatus(node, statusId, checkId string) (status, output string)
	CheckKeyExists(key string) bool

	GetProfileInfo(node, serviceID, checkID string) ProfileInfo

	GetReminders() []notifier.Message
	SetReminder(m notifier.Message)
	DeleteReminder(node string, checkid string)
}

// DefaultAlertConfig loads default config settings
func DefaultAlertConfig() *ConsulAlertConfig {

	checks := &ChecksConfig{
		Enabled:         true,
		ChangeThreshold: 60,
	}

	events := &EventsConfig{
		Enabled:  true,
		Handlers: []string{},
	}

	email := &notifier.EmailNotifier{
		ClusterName: "Consul-Alerts",
		Enabled:     false,
		SenderAlias: "Consul Alerts",
		Receivers:   []string{},
	}

	log := &notifier.LogNotifier{
		Enabled: true,
		Path:    "/tmp/consul-notifications.log",
	}

	influxdb := &notifier.InfluxdbNotifier{
		Enabled:    false,
		SeriesName: "consul-alerts",
	}

	slack := &notifier.SlackNotifier{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	mattermost := &notifier.MattermostNotifier{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	pagerduty := &notifier.PagerDutyNotifier{
		Enabled: false,
	}

	hipchat := &notifier.HipChatNotifier{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	opsgenie := &notifier.OpsGenieNotifier{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	telegram := &notifier.TelegramNotifier{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	awsSns := &notifier.AwsSnsNotifier{
		Enabled: false,
	}

	victorOps := &notifier.VictorOpsNotifier{
		Enabled: false,
	}

	notifiers := &notifier.Notifiers{
		Email:      email,
		Log:        log,
		Influxdb:   influxdb,
		Slack:      slack,
		Mattermost: mattermost,
		PagerDuty:  pagerduty,
		HipChat:    hipchat,
		OpsGenie:   opsgenie,
		Telegram:   telegram,
		AwsSns:     awsSns,
		VictorOps:  victorOps,
		Custom:     []string{},
	}

	return &ConsulAlertConfig{
		Checks:    checks,
		Events:    events,
		Notifiers: notifiers,
	}
}

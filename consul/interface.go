package consul

import (
	"time"

	"github.com/AcalephStorage/consul-alerts/notifier"
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
	Notifiers *NotifiersConfig
}

type ChecksConfig struct {
	Enabled         bool
	ChangeThreshold int
}

type EventsConfig struct {
	Enabled  bool
	Handlers []string
}

type NotifiersConfig struct {
	Email     *EmailNotifierConfig
	Log       *LogNotifierConfig
	Influxdb  *InfluxdbNotifierConfig
	Slack     *SlackNotifierConfig
	PagerDuty *PagerDutyNotifierConfig
	HipChat   *HipChatNotifierConfig
	OpsGenie  *OpsGenieNotifierConfig
	Custom    []string
}

type EmailNotifierConfig struct {
	ClusterName string
	Enabled     bool
	Url         string
	Port        int
	Username    string
	Password    string
	SenderAlias string
	SenderEmail string
	Receivers   []string
	Template    string
}

type LogNotifierConfig struct {
	Enabled bool
	Path    string
}

type InfluxdbNotifierConfig struct {
	Enabled    bool
	Host       string
	Username   string
	Password   string
	Database   string
	SeriesName string
}

type SlackNotifierConfig struct {
	Enabled     bool
	ClusterName string
	Url         string
	Channel     string
	Username    string
	IconUrl     string
	IconEmoji   string
	Detailed    bool
}

type PagerDutyNotifierConfig struct {
	Enabled    bool
	ServiceKey string
	ClientName string
	ClientUrl  string
}

type HipChatNotifierConfig struct {
	Enabled     bool
	ClusterName string
	RoomId      string
	AuthToken   string
	BaseURL     string
	From        string
}

type OpsGenieNotifierConfig struct {
	Enabled     bool
	ClusterName string
	ApiKey      string
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
	Interval  int
	NotifList map[string]bool
}

// Consul interface provides access to consul client 
type Consul interface {
	LoadConfig()

	EventsEnabled() bool
	ChecksEnabled() bool
	EventHandlers(eventName string) []string

	EmailConfig() *EmailNotifierConfig
	LogConfig() *LogNotifierConfig
	InfluxdbConfig() *InfluxdbNotifierConfig
	SlackConfig() *SlackNotifierConfig
	PagerDutyConfig() *PagerDutyNotifierConfig
	HipChatConfig() *HipChatNotifierConfig
	OpsGenieConfig() *OpsGenieNotifierConfig

	CheckChangeThreshold() int
	UpdateCheckData()
	NewAlerts() []Check

	IsBlacklisted(check *Check) bool

	CustomNotifiers() map[string]string

	CheckStatus(node, statusId, checkId string) (status, output string)
	CheckKeyExists(key string) bool

	GetProfileInfo(node, serviceID, checkID string) (notifiersList map[string]bool, interval int)
	
	GetReminders() []notifier.Message
	SetReminder(m notifier.Message)
	DeleteReminder(node string)
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

	email := &EmailNotifierConfig{
		ClusterName: "Consul-Alerts",
		Enabled:     false,
		SenderAlias: "Consul Alerts",
		Receivers:   []string{},
	}

	log := &LogNotifierConfig{
		Enabled: true,
		Path:    "/tmp/consul-notifications.log",
	}

	influxdb := &InfluxdbNotifierConfig{
		Enabled:    false,
		SeriesName: "consul-alerts",
	}

	slack := &SlackNotifierConfig{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	pagerduty := &PagerDutyNotifierConfig{
		Enabled: false,
	}

	hipchat := &HipChatNotifierConfig{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	opsgenie := &OpsGenieNotifierConfig{
		Enabled:     false,
		ClusterName: "Consul-Alerts",
	}

	notifiers := &NotifiersConfig{
		Email:     email,
		Log:       log,
		Influxdb:  influxdb,
		Slack:     slack,
		PagerDuty: pagerduty,
		HipChat:   hipchat,
		OpsGenie:  opsgenie,
		Custom:    []string{},
	}

	return &ConsulAlertConfig{
		Checks:    checks,
		Events:    events,
		Notifiers: notifiers,
	}
}

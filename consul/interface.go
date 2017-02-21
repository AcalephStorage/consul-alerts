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

type NotifiersConfig struct {
	Email     *EmailNotifierConfig
	Log       *LogNotifierConfig
	Influxdb  *InfluxdbNotifierConfig
	Slack     *SlackNotifierConfig
	PagerDuty *PagerDutyNotifierConfig
	HipChat   *HipChatNotifierConfig
	OpsGenie  *OpsGenieNotifierConfig
	AwsSns    *AwsSnsNotifierConfig
	VictorOps *VictorOpsNotifierConfig
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
	Receivers   map[string][]string
	Template    string
	OnePerAlert bool
	OnePerNode  bool
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

type AwsSnsNotifierConfig struct {
	Enabled  bool
	Region   string
	TopicArn string
}

// VictorOpsNotifierConfig provides configuration options for VictorOps notifier
type VictorOpsNotifierConfig struct {
	Enabled    bool
	APIKey     string
	RoutingKey string
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
	Interval      int
	NotifList     map[string]bool
	NotifTypeList map[string][]string
	VarOverrides  notifier.Notifiers
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
	PagerDutyNotifier() *notifier.PagerDutyNotifier
	HipChatNotifier() *notifier.HipChatNotifier
	OpsGenieNotifier() *notifier.OpsGenieNotifier
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
		Receivers:   map[string][]string{},
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

	awsSns := &notifier.AwsSnsNotifier{
		Enabled: false,
	}

	victorOps := &notifier.VictorOpsNotifier{
		Enabled: false,
	}

	notifiers := &notifier.Notifiers{
		Email:     email,
		Log:       log,
		Influxdb:  influxdb,
		Slack:     slack,
		PagerDuty: pagerduty,
		HipChat:   hipchat,
		OpsGenie:  opsgenie,
		AwsSns:    awsSns,
		VictorOps: victorOps,
		Custom:    []string{},
	}

	return &ConsulAlertConfig{
		Checks:    checks,
		Events:    events,
		Notifiers: notifiers,
	}
}

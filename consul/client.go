package consul

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"regexp"

	"encoding/json"

	"github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	consulapi "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/hashicorp/consul/api"
)

const (
	ConfigTypeBool = iota
	ConfigTypeString
	ConfigTypeInt
	ConfigTypeStrArray
)

type configType int

type ConsulAlertClient struct {
	api    *consulapi.Client
	config *ConsulAlertConfig
}

func NewClient(address, dc, aclToken string) (*ConsulAlertClient, error) {
	config := consulapi.DefaultConfig()
	config.Address = address
	config.Datacenter = dc
	config.Token = aclToken
	api, _ := consulapi.NewClient(config)
	alertConfig := DefaultAlertConfig()

	client := &ConsulAlertClient{
		api:    api,
		config: alertConfig,
	}

	try := 1
	for {
		try += try
		log.Println("Checking consul agent connection...")
		_, err := client.api.Status().Leader()
		if err != nil {
			log.Println("Waiting for consul:", err)
			if try > 10 {
				return nil, err
			}
			time.Sleep(10000 * time.Millisecond)
		} else {
			break
		}
	}

	client.LoadConfig()
	client.UpdateCheckData()
	return client, nil
}

func (c *ConsulAlertClient) LoadConfig() {
	if kvPairs, _, err := c.api.KV().List("consul-alerts/config", nil); err == nil {

		config := c.config

		for _, kvPair := range kvPairs {

			key := kvPair.Key
			val := kvPair.Value

			var valErr error
			switch key {
			// checks config
			case "consul-alerts/config/checks/enabled":
				valErr = loadCustomValue(&config.Checks.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/checks/change-threshold":
				valErr = loadCustomValue(&config.Checks.ChangeThreshold, val, ConfigTypeInt)

			// events config
			case "consul-alerts/config/events/enabled":
				valErr = loadCustomValue(&config.Events.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/events/handlers":
				valErr = loadCustomValue(&config.Events.Handlers, val, ConfigTypeStrArray)

			// email notifier config
			case "consul-alerts/config/notifiers/email/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.Email.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/template":
				valErr = loadCustomValue(&config.Notifiers.Email.Template, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/enabled":
				valErr = loadCustomValue(&config.Notifiers.Email.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/email/password":
				valErr = loadCustomValue(&config.Notifiers.Email.Password, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/port":
				valErr = loadCustomValue(&config.Notifiers.Email.Port, val, ConfigTypeInt)
			case "consul-alerts/config/notifiers/email/receivers":
				valErr = loadCustomValue(&config.Notifiers.Email.Receivers, val, ConfigTypeStrArray)
			case "consul-alerts/config/notifiers/email/sender-alias":
				valErr = loadCustomValue(&config.Notifiers.Email.SenderAlias, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/sender-email":
				valErr = loadCustomValue(&config.Notifiers.Email.SenderEmail, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/url":
				valErr = loadCustomValue(&config.Notifiers.Email.Url, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/username":
				valErr = loadCustomValue(&config.Notifiers.Email.Username, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/email/one-per-alert":
				valErr = loadCustomValue(&config.Notifiers.Email.OnePerAlert, val, ConfigTypeBool)

			// log notifier config
			case "consul-alerts/config/notifiers/log/enabled":
				valErr = loadCustomValue(&config.Notifiers.Log.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/log/path":
				valErr = loadCustomValue(&config.Notifiers.Log.Path, val, ConfigTypeString)

			// influxdb notifier config
			case "consul-alerts/config/notifiers/influxdb/enabled":
				valErr = loadCustomValue(&config.Notifiers.Influxdb.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/influxdb/host":
				valErr = loadCustomValue(&config.Notifiers.Influxdb.Host, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/influxdb/username":
				valErr = loadCustomValue(&config.Notifiers.Influxdb.Username, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/influxdb/password":
				valErr = loadCustomValue(&config.Notifiers.Influxdb.Password, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/influxdb/database":
				valErr = loadCustomValue(&config.Notifiers.Influxdb.Database, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/influxdb/series-name":
				valErr = loadCustomValue(&config.Notifiers.Influxdb.SeriesName, val, ConfigTypeString)

			// slack notfier config
			case "consul-alerts/config/notifiers/slack/enabled":
				valErr = loadCustomValue(&config.Notifiers.Slack.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/slack/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.Slack.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/slack/url":
				valErr = loadCustomValue(&config.Notifiers.Slack.Url, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/slack/channel":
				valErr = loadCustomValue(&config.Notifiers.Slack.Channel, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/slack/username":
				valErr = loadCustomValue(&config.Notifiers.Slack.Username, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/slack/icon-url":
				valErr = loadCustomValue(&config.Notifiers.Slack.IconUrl, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/slack/icon-emoji":
				valErr = loadCustomValue(&config.Notifiers.Slack.IconEmoji, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/slack/detailed":
				valErr = loadCustomValue(&config.Notifiers.Slack.Detailed, val, ConfigTypeBool)

			// pager-duty notfier config
			case "consul-alerts/config/notifiers/pagerduty/enabled":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/pagerduty/service-key":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.ServiceKey, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/pagerduty/client-name":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.ClientName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/pagerduty/client-url":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.ClientUrl, val, ConfigTypeString)

			// hipchat notfier config
			case "consul-alerts/config/notifiers/hipchat/enabled":
				valErr = loadCustomValue(&config.Notifiers.HipChat.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/hipchat/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.HipChat.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/hipchat/room-id":
				valErr = loadCustomValue(&config.Notifiers.HipChat.RoomId, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/hipchat/auth-token":
				valErr = loadCustomValue(&config.Notifiers.HipChat.AuthToken, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/hipchat/base-url":
				valErr = loadCustomValue(&config.Notifiers.HipChat.BaseURL, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/hipchat/from":
				valErr = loadCustomValue(&config.Notifiers.HipChat.From, val, ConfigTypeString)

			// OpsGenie notifier config
			case "consul-alerts/config/notifiers/opsgenie/enabled":
				valErr = loadCustomValue(&config.Notifiers.OpsGenie.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/opsgenie/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.OpsGenie.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/opsgenie/api-key":
				valErr = loadCustomValue(&config.Notifiers.OpsGenie.ApiKey, val, ConfigTypeString)

			// AwsSns notifier config
			case "consul-alerts/config/notifiers/awssns/enabled":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/awssns/region":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.Region, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/awssns/topic-arn":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.TopicArn, val, ConfigTypeString)

			}

			if valErr != nil {
				log.Printf(`unable to load custom value for "%s". Using default instead. Error: %s`, key, valErr.Error())
			}

		}
	} else {
		log.Println("Unable to load custom config, using default instead:", err)
	}

}

func loadCustomValue(configVariable interface{}, data []byte, cType configType) (err error) {
	switch cType {
	case ConfigTypeBool:
		var val bool
		if val, err = strconv.ParseBool(string(data)); err == nil {
			boolConfig := configVariable.(*bool)
			*boolConfig = val
		}
	case ConfigTypeString:
		strConfig := configVariable.(*string)
		*strConfig = string(data)
	case ConfigTypeInt:
		var val int
		if val, err = strconv.Atoi(string(data)); err == nil {
			intConfig := configVariable.(*int)
			*intConfig = int(val)
		}
	case ConfigTypeStrArray:
		arrConfig := configVariable.(*[]string)
		err = json.Unmarshal(data, arrConfig)
	}
	return err
}

func (c *ConsulAlertClient) EventsEnabled() bool {
	return c.config.Events.Enabled
}

func (c *ConsulAlertClient) ChecksEnabled() bool {
	return c.config.Checks.Enabled
}

func (c *ConsulAlertClient) EventHandlers(eventName string) []string {
	return c.config.Events.Handlers
}

func (c *ConsulAlertClient) CheckChangeThreshold() int {
	return c.config.Checks.ChangeThreshold
}

func (c *ConsulAlertClient) UpdateCheckData() {
	healthApi := c.api.Health()
	kvApi := c.api.KV()

	healths, _, _ := healthApi.State("any", nil)

	for _, health := range healths {

		node := health.Node
		service := health.ServiceID
		check := health.CheckID
		if service == "" {
			service = "_"
		}
		key := fmt.Sprintf("consul-alerts/checks/%s/%s/%s", node, service, check)

		status, _, _ := kvApi.Get(key, nil)
		existing := status != nil

		localHealth := Check(*health)

		if c.IsBlacklisted(&localHealth) {
			log.Printf("%s:%s:%s is blacklisted.", node, service, check)
			continue
		}

		if !existing {
			c.registerHealthCheck(key, &localHealth)
		} else {
			c.updateHealthCheck(key, &localHealth)
		}

	}

}

// GetReminders returns list of reminders
func (c *ConsulAlertClient) GetReminders() []notifier.Message {
	remindersList, _, _ := c.api.KV().List("consul-alerts/reminders", nil)
	var messages []notifier.Message
	for _, kvpair := range remindersList {
		var message notifier.Message
		json.Unmarshal(kvpair.Value, &message)
		messages = append(messages, message)
	}
	log.Println("Getting reminders")
	return messages
}

// SetReminder sets a reminder
func (c *ConsulAlertClient) SetReminder(m notifier.Message) {
	data, _ := json.Marshal(m)
	key := fmt.Sprintf("consul-alerts/reminders/%s", m.Node)
	c.api.KV().Put(&consulapi.KVPair{Key: key, Value: data}, nil)
	log.Println("Setting reminder for node: ", m.Node)
}

// DeleteReminder deletes a reminder
func (c *ConsulAlertClient) DeleteReminder(node string) {
	key := fmt.Sprintf("consul-alerts/reminders/%s", node)
	c.api.KV().Delete(key, nil)
	log.Println("Deleting reminder for node: ", node)
}

// NewAlerts returns a list of checks marked for notification
func (c *ConsulAlertClient) NewAlerts() []Check {
	allChecks, _, _ := c.api.KV().List("consul-alerts/checks", nil)
	var alerts []Check
	for _, kvpair := range allChecks {
		key := kvpair.Key
		if strings.HasSuffix(key, "/") {
			continue
		}
		var status Status
		json.Unmarshal(kvpair.Value, &status)
		if status.ForNotification {
			status.ForNotification = false
			data, _ := json.Marshal(status)
			c.api.KV().Put(&consulapi.KVPair{Key: key, Value: data}, nil)
			// check if blacklisted

			if !c.IsBlacklisted(status.HealthCheck) {
				alerts = append(alerts, *status.HealthCheck)
			}
		}
	}
	return alerts
}

// CustomNotifiers returns a map of all custom notifiers and command path as the key value
func (c *ConsulAlertClient) CustomNotifiers() (customNotifs map[string]string) {
	if kvPairs, _, err := c.api.KV().List("consul-alerts/config/notifiers/custom", nil); err == nil {
		customNotifs = make(map[string]string)
		for _, kvPair := range kvPairs {
			rp := regexp.MustCompile("/([^/]*)$")
			match := rp.FindStringSubmatch(kvPair.Key)
			custNotifName := match[1]
			if custNotifName == "" {
				continue
			}
			customNotifs[custNotifName] = string(kvPair.Value)
		}
	}
	return customNotifs
}

// EmailConfig exports the email config
func (c *ConsulAlertClient) EmailConfig() *EmailNotifierConfig {
	return c.config.Notifiers.Email
}

func (c *ConsulAlertClient) LogConfig() *LogNotifierConfig {
	return c.config.Notifiers.Log
}

func (c *ConsulAlertClient) InfluxdbConfig() *InfluxdbNotifierConfig {
	return c.config.Notifiers.Influxdb
}

func (c *ConsulAlertClient) SlackConfig() *SlackNotifierConfig {
	return c.config.Notifiers.Slack
}

func (c *ConsulAlertClient) PagerDutyConfig() *PagerDutyNotifierConfig {
	return c.config.Notifiers.PagerDuty
}

func (c *ConsulAlertClient) HipChatConfig() *HipChatNotifierConfig {
	return c.config.Notifiers.HipChat
}

func (c *ConsulAlertClient) OpsGenieConfig() *OpsGenieNotifierConfig {
	return c.config.Notifiers.OpsGenie
}

func (c *ConsulAlertClient) AwsSnsConfig() *AwsSnsNotifierConfig {
	return c.config.Notifiers.AwsSns
}

func (c *ConsulAlertClient) registerHealthCheck(key string, health *Check) {

	log.Printf(
		"Registering new health check: node=%s, service=%s, check=%s, status=%s",
		health.Node,
		health.ServiceName,
		health.Name,
		health.Status,
	)

	var newStatus Status
	if health.Status == "passing" {
		newStatus = Status{
			Current:          health.Status,
			CurrentTimestamp: time.Now(),
			HealthCheck:      health,
		}
	} else {
		newStatus = Status{
			Pending:          health.Status,
			PendingTimestamp: time.Now(),
			HealthCheck:      health,
		}
	}

	statusData, _ := json.Marshal(newStatus)
	c.api.KV().Put(&consulapi.KVPair{Key: key, Value: statusData}, nil)
}

func (c *ConsulAlertClient) updateHealthCheck(key string, health *Check) {

	kvpair, _, _ := c.api.KV().Get(key, nil)
	val := kvpair.Value
	var storedStatus Status
	json.Unmarshal(val, &storedStatus)

	// no status change if the stored status and latest status is the same
	noStatusChange := storedStatus.Current == health.Status

	// new pending status if it's a new status and it's not the same as the pending status
	newPendingStatus := storedStatus.Current != health.Status && storedStatus.Pending != health.Status

	// status is still pending for change. will change if it reaches threshold
	stillPendingStatus := storedStatus.Current != health.Status && storedStatus.Pending == health.Status

	switch {

	case noStatusChange:
		if storedStatus.Pending != "" {
			storedStatus.Pending = ""
			storedStatus.PendingTimestamp = time.Time{}
			log.Printf(
				"%s:%s:%s is now back to %s.",
				health.Node,
				health.ServiceName,
				health.Name,
				storedStatus.Current,
			)
		}

	case newPendingStatus:
		storedStatus.Pending = health.Status
		storedStatus.PendingTimestamp = time.Now()
		log.Printf(
			"%s:%s:%s is now pending status change from %s to %s.",
			health.Node,
			health.ServiceName,
			health.Name,
			storedStatus.Current,
			storedStatus.Pending,
		)

	case stillPendingStatus:
		duration := time.Since(storedStatus.PendingTimestamp)
		if int(duration.Seconds()) >= c.config.Checks.ChangeThreshold {

			log.Printf(
				"%s:%s:%s has changed status from %s to %s.",
				health.Node,
				health.ServiceName,
				health.Name,
				storedStatus.Current,
				storedStatus.Pending,
			)

			storedStatus.Current = storedStatus.Pending
			storedStatus.CurrentTimestamp = time.Now()
			storedStatus.Pending = ""
			storedStatus.PendingTimestamp = time.Time{}
			storedStatus.ForNotification = true
		} else {
			log.Printf(
				"%s:%s:%s is pending status change from %s to %s for %s.",
				health.Node,
				health.ServiceName,
				health.Name,
				storedStatus.Current,
				storedStatus.Pending,
				duration,
			)
		}

	}
	storedStatus.HealthCheck = health

	data, _ := json.Marshal(storedStatus)
	c.api.KV().Put(&consulapi.KVPair{Key: key, Value: data}, nil)
}

func (c *ConsulAlertClient) CheckStatus(node, serviceId, checkId string) (status, output string) {
	if serviceId == "" {
		serviceId = "_"
	}
	key := fmt.Sprintf("consul-alerts/checks/%s/%s/%s", node, serviceId, checkId)
	kvPair, _, _ := c.api.KV().Get(key, nil)

	if kvPair == nil {
		status = ""
		output = ""
		return
	}

	var checkStatus Status
	json.Unmarshal(kvPair.Value, &checkStatus)

	status = checkStatus.Current
	output = checkStatus.HealthCheck.Output
	return
}

// GetProfileInfo returns profile info for check
func (c *ConsulAlertClient) GetProfileInfo(node, serviceID, checkID string) (notifiersList map[string]bool, interval int) {
	log.Println("Getting profile for node: ", node, " service: ", serviceID, " check: ", checkID)

	var profile string

	kvPair, _, _ := c.api.KV().Get(fmt.Sprintf("consul-alerts/config/notif-selection/services/%s", serviceID), nil)
	if kvPair != nil {
		profile = string(kvPair.Value)
		log.Println("service selection key found.")
	} else if kvPair, _, _ = c.api.KV().Get(fmt.Sprintf("consul-alerts/config/notif-selection/checks/%s", checkID), nil); kvPair != nil {
		profile = string(kvPair.Value)
		log.Println("check selection key found.")
	} else if kvPair, _, _ = c.api.KV().Get(fmt.Sprintf("consul-alerts/config/notif-selection/hosts/%s", node), nil); kvPair != nil {
		profile = string(kvPair.Value)
		log.Println("host selection key found.")
	} else {
		profile = "default"
	}

	key := fmt.Sprintf("consul-alerts/config/notif-profiles/%s", profile)
	log.Println("profile key: ", key)
	kvPair, _, _ = c.api.KV().Get(key, nil)
	if kvPair == nil {
		log.Println("profile key not found.")
		return
	}
	var checkProfile ProfileInfo
	json.Unmarshal(kvPair.Value, &checkProfile)

	notifiersList = checkProfile.NotifList
	interval = checkProfile.Interval
	log.Println("Interval: ", interval, " List: ", notifiersList)
	return
}

// IsBlacklisted gets the blacklist status of check
func (c *ConsulAlertClient) IsBlacklisted(check *Check) bool {
	node := check.Node
	nodeCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/nodes/%s", node)
	nodeBlacklisted := c.CheckKeyExists(nodeCheckKey)

	service := "_"
	serviceBlacklisted := false
	if check.ServiceID != "" {
		service = check.ServiceID
		serviceCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/services/%s", service)
		serviceBlacklisted = c.CheckKeyExists(serviceCheckKey)
	}

	checkId := check.CheckID
	checkCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/checks/%s", checkId)
	checkBlacklisted := c.CheckKeyExists(checkCheckKey)

	singleKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/single/%s/%s/%s", node, service, checkId)
	singleBlacklisted := c.CheckKeyExists(singleKey)

	return nodeBlacklisted || serviceBlacklisted || checkBlacklisted || singleBlacklisted
}

func (c *ConsulAlertClient) CheckKeyExists(key string) bool {
	kvpair, _, err := c.api.KV().Get(key, nil)
	return kvpair != nil && err == nil
}

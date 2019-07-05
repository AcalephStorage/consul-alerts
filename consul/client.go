package consul

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"regexp"

	"encoding/json"

	notifier "github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	consulapi "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/hashicorp/consul/api"
)

const (
	ConfigTypeBool = iota
	ConfigTypeString
	ConfigTypeInt
	ConfigTypeStrArray
	ConfigTypeStrMap
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
			case "consul-alerts/config/notifiers/email/receivers/":
				kvmTemp := c.KvMap("consul-alerts/config/notifiers/email/receivers")
				// only want the key at the end, so split on slashes and take the last item
				kvm := make(map[string][]string, len(kvmTemp))
				for k, v := range kvmTemp {
					kSplit := strings.Split(k, "/")
					kvm[kSplit[len(kSplit)-1]] = v
				}
				convertedVal, _ := json.Marshal(kvm)
				valErr = loadCustomValue(&config.Notifiers.Email.Receivers, convertedVal, ConfigTypeStrMap)
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
			case "consul-alerts/config/notifiers/email/one-per-node":
				valErr = loadCustomValue(&config.Notifiers.Email.OnePerNode, val, ConfigTypeBool)

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

				// mattermost notfier config
			case "consul-alerts/config/notifiers/mattermost/enabled":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/mattermost/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost/url":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.Url, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost/username":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.UserName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost/password":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.Password, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost/team":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.Team, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost/channel":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.Channel, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost/detailed":
				valErr = loadCustomValue(&config.Notifiers.Mattermost.Detailed, val, ConfigTypeBool)

				// mattermost webhook notifier config
			case "consul-alerts/config/notifiers/mattermost-webhook/enabled":
				valErr = loadCustomValue(&config.Notifiers.MattermostWebhook.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/mattermost-webhook/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.MattermostWebhook.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost-webhook/url":
				valErr = loadCustomValue(&config.Notifiers.MattermostWebhook.Url, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost-webhook/channel":
				valErr = loadCustomValue(&config.Notifiers.MattermostWebhook.Channel, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost-webhook/username":
				valErr = loadCustomValue(&config.Notifiers.MattermostWebhook.Username, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/mattermost-webhook/icon-url":
				valErr = loadCustomValue(&config.Notifiers.MattermostWebhook.IconUrl, val, ConfigTypeString)

				// pager-duty notfier config
			case "consul-alerts/config/notifiers/pagerduty/enabled":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/pagerduty/service-key":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.ServiceKey, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/pagerduty/client-name":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.ClientName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/pagerduty/client-url":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.ClientUrl, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/pagerduty/max-retry":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.MaxRetry, val, ConfigTypeInt)
			case "consul-alerts/config/notifiers/pagerduty/retry-base-interval":
				valErr = loadCustomValue(&config.Notifiers.PagerDuty.RetryBaseInterval, val, ConfigTypeInt)

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
			case "consul-alerts/config/notifiers/awssns/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/awssns/enabled":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/awssns/region":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.Region, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/awssns/topic-arn":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.TopicArn, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/awssns/template":
				valErr = loadCustomValue(&config.Notifiers.AwsSns.Template, val, ConfigTypeString)

				// VictorOps notfier config
			case "consul-alerts/config/notifiers/victorops/enabled":
				valErr = loadCustomValue(&config.Notifiers.VictorOps.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/victorops/api-key":
				valErr = loadCustomValue(&config.Notifiers.VictorOps.APIKey, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/victorops/routing-key":
				valErr = loadCustomValue(&config.Notifiers.VictorOps.RoutingKey, val, ConfigTypeString)

			// http endpoint notfier config
			case "consul-alerts/config/notifiers/http-endpoint/enabled":
				valErr = loadCustomValue(&config.Notifiers.HttpEndpoint.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/http-endpoint/cluster-name":
				valErr = loadCustomValue(&config.Notifiers.HttpEndpoint.ClusterName, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/http-endpoint/base-url":
				valErr = loadCustomValue(&config.Notifiers.HttpEndpoint.BaseURL, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/http-endpoint/endpoint":
				valErr = loadCustomValue(&config.Notifiers.HttpEndpoint.Endpoint, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/http-endpoint/payload":
				valErr = loadCustomValue(&config.Notifiers.HttpEndpoint.Payload, val, ConfigTypeStrMap)

			// iLert notfier config
			case "consul-alerts/config/notifiers/ilert/enabled":
				valErr = loadCustomValue(&config.Notifiers.ILert.Enabled, val, ConfigTypeBool)
			case "consul-alerts/config/notifiers/ilert/api-key":
				valErr = loadCustomValue(&config.Notifiers.ILert.ApiKey, val, ConfigTypeString)
			case "consul-alerts/config/notifiers/ilert/incident-key-template":
				valErr = loadCustomValue(&config.Notifiers.ILert.IncidentKeyTemplate, val, ConfigTypeString)
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
		if val, err = strconv.ParseBool(strings.Trim(string(data), " \t\n")); err == nil {
			boolConfig := configVariable.(*bool)
			*boolConfig = val
		}
	case ConfigTypeString:
		strConfig := configVariable.(*string)
		*strConfig = strings.Trim(string(data), " \t\n")
	case ConfigTypeInt:
		var val int
		if val, err = strconv.Atoi(strings.Trim(string(data), " \t\n")); err == nil {
			intConfig := configVariable.(*int)
			*intConfig = int(val)
		}
	case ConfigTypeStrArray:
		arrConfig := configVariable.(*[]string)
		err = json.Unmarshal(data, arrConfig)
	case ConfigTypeStrMap:
		mapConfig := configVariable.(*map[string][]string)
		err = json.Unmarshal(data, mapConfig)
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
	reminderKeys, _, _ := c.api.KV().List("consul-alerts/reminders/", nil)
	remindersSubsLevel := 4

	for index := range reminderKeys {
		log.Printf("checking for stale reminders")
		s := strings.Split(reminderKeys[index].Key, "/")
		// check if the consul-alerts/reminders/ folder has sub folders
		if len(s) >= remindersSubsLevel {
			node, check := s[2], s[3]

			nodecat, _, _ := c.api.Health().Node(node, nil)
			settodelete := true

			for j := range nodecat {
				if nodecat[j].CheckID == check {
					settodelete = false
					break
				}
			}
			if settodelete {
				log.Printf("Reminder %s %s needs to be deleted, stale", node, check)
				c.DeleteReminder(node, check)
			}
		}
	}

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

		reminderkey := fmt.Sprintf("consul-alerts/reminders/%s/%s", node, check)
		reminderstatus, _, err := kvApi.Get(reminderkey, nil)
		reminderexists := reminderstatus != nil

		if err != nil {
			log.Println("Unable to get kv value: ", err)
		}

		if reminderexists {

			var remindermap map[string]interface{}

			json.Unmarshal((reminderstatus.Value), &remindermap)

			if remindermap["Output"] != health.Output {
				log.Printf("Updating reminder data for %s", reminderkey)

				remindermap["Output"] = health.Output
				newreminder, _ := json.Marshal(remindermap)
				reminderstatus.Value = newreminder

				_, _, err := kvApi.CAS(reminderstatus, nil)
				if err != nil {
					log.Println("Unable to set kv value: ", err)
				}
			}
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
	key := fmt.Sprintf("consul-alerts/reminders/%s/%s", m.Node, m.CheckId)
	c.api.KV().Put(&consulapi.KVPair{Key: key, Value: data}, nil)
	log.Println("Setting reminder for node: ", m.Node)
}

// DeleteReminder deletes a reminder
func (c *ConsulAlertClient) DeleteReminder(node string, checkid string) {
	key := fmt.Sprintf("consul-alerts/reminders/%s/%s", node, checkid)
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

// KvMap returns a map of KV pairs found directly inside the passed path
func (c *ConsulAlertClient) KvMap(kvPath string) (kvMap map[string][]string) {
	if kvPairs, _, err := c.api.KV().List(kvPath, nil); err == nil {
		kvMap = make(map[string][]string)
		for _, kvPair := range kvPairs {
			if strings.HasSuffix(kvPair.Key, "/") {
				continue
			}
			itemList := []string{}
			json.Unmarshal(kvPair.Value, &itemList)
			kvMap[string(kvPair.Key)] = itemList
		}
	}
	return kvMap
}

func (c *ConsulAlertClient) NewAlertsWithFilter(nodeName string, serviceName string, checkName string, statuses []string, ignoreBlacklist bool) []Check {
	allChecks, _, _ := c.api.KV().List("consul-alerts/checks", nil)
	alerts := make([]Check, 0)
	for _, kvpair := range allChecks {
		if strings.HasSuffix(kvpair.Key, "/") {
			continue
		}

		var status Status
		json.Unmarshal(kvpair.Value, &status)

		check := *status.HealthCheck

		if nodeName != "" && nodeName != check.Node {
			continue
		}

		if serviceName != "" && serviceName != check.ServiceName {
			continue
		}

		if checkName != "" && checkName != check.Name {
			continue
		}

		if len(statuses) > 0 {
			inStatuses := false
			for _, s := range statuses {
				inStatuses = check.Status == s
			}
			if !inStatuses {
				continue
			}
		}

		if !ignoreBlacklist && c.IsBlacklisted(status.HealthCheck) {
			continue
		}
		alerts = append(alerts, *status.HealthCheck)
	}
	return alerts
}

func (c *ConsulAlertClient) EmailNotifier() *notifier.EmailNotifier {
	return c.config.Notifiers.Email
}

func (c *ConsulAlertClient) LogNotifier() *notifier.LogNotifier {
	return c.config.Notifiers.Log
}

func (c *ConsulAlertClient) InfluxdbNotifier() *notifier.InfluxdbNotifier {
	return c.config.Notifiers.Influxdb
}

func (c *ConsulAlertClient) SlackNotifier() *notifier.SlackNotifier {
	return c.config.Notifiers.Slack
}

func (c *ConsulAlertClient) MattermostNotifier() *notifier.MattermostNotifier {
	return c.config.Notifiers.Mattermost
}

func (c *ConsulAlertClient) MattermostWebhookNotifier() *notifier.MattermostWebhookNotifier {
	return c.config.Notifiers.MattermostWebhook
}

func (c *ConsulAlertClient) PagerDutyNotifier() *notifier.PagerDutyNotifier {
	return c.config.Notifiers.PagerDuty
}

func (c *ConsulAlertClient) HipChatNotifier() *notifier.HipChatNotifier {
	return c.config.Notifiers.HipChat
}

func (c *ConsulAlertClient) OpsGenieNotifier() *notifier.OpsGenieNotifier {
	return c.config.Notifiers.OpsGenie
}

func (c *ConsulAlertClient) AwsSnsNotifier() *notifier.AwsSnsNotifier {
	return c.config.Notifiers.AwsSns
}

func (c *ConsulAlertClient) VictorOpsNotifier() *notifier.VictorOpsNotifier {
	return c.config.Notifiers.VictorOps
}

func (c *ConsulAlertClient) HttpEndpointNotifier() *notifier.HttpEndpointNotifier {
	return c.config.Notifiers.HttpEndpoint
}

func (c *ConsulAlertClient) ILertNotifier() *notifier.ILertNotifier {
	return c.config.Notifiers.ILert
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

	// indicate whether we are changing storedStatus to prevent unnecessary PUT to KV
	changed := false

	switch {

	case noStatusChange:
		if storedStatus.Pending != "" {
			storedStatus.Pending = ""
			storedStatus.PendingTimestamp = time.Time{}
			changed = true
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
		changed = true
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

		changeThreshold := c.config.Checks.ChangeThreshold
		if override := c.GetChangeThreshold(health); override >= 0 {
			changeThreshold = override
		}

		if int(duration.Seconds()) >= changeThreshold {

			log.Printf(
				"%s:%s:%s has changed status from %s to %s.",
				health.Node,
				health.ServiceName,
				health.Name,
				storedStatus.Current,
				storedStatus.Pending,
			)

			// do not trigger a notification if the check was just registered and
			// the first status is passing
			if storedStatus.Current == "" && storedStatus.Pending == "passing" {
				storedStatus.ForNotification = false
			} else {
				storedStatus.ForNotification = true
			}

			storedStatus.Current = storedStatus.Pending
			storedStatus.CurrentTimestamp = time.Now()
			storedStatus.Pending = ""
			storedStatus.PendingTimestamp = time.Time{}
			changed = true
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

	if !changed {
		return
	}

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

// getProfileForEntity returns the profile matching the exact path or the regexp
// entity is either 'service', 'check' or 'host'
func (c *ConsulAlertClient) getProfileForEntity(entity string, id string) string {
	kvPair, _, _ := c.api.KV().Get(
		fmt.Sprintf("consul-alerts/config/notif-selection/%ss/%s",
			entity, id), nil)
	if kvPair != nil {
		log.Printf("%s selection key found.\n", entity)
		return string(kvPair.Value)
	} else if kvPair, _, _ := c.api.KV().Get(
		fmt.Sprintf("consul-alerts/config/notif-selection/%ss", entity),
		nil); kvPair != nil {
		var regexMap map[string]string
		json.Unmarshal(kvPair.Value, &regexMap)
		for pattern, profile := range regexMap {
			matched, err := regexp.MatchString(pattern, id)
			if err != nil {
				log.Printf("unable to match %s %s against pattern %s. Error: %s\n",
					entity, id, pattern, err.Error())
			} else if matched {
				log.Printf("Regexp matching %s found (%s).\n", entity, pattern)
				return profile
			}
		}
	}
	return ""
}

func (c *ConsulAlertClient) getProfileForService(serviceID string) string {
	return c.getProfileForEntity("service", serviceID)
}

func (c *ConsulAlertClient) getProfileForCheck(checkID string) string {
	return c.getProfileForEntity("check", checkID)
}

func (c *ConsulAlertClient) getProfileForNode(node string) string {
	return c.getProfileForEntity("host", node)
}

func (c *ConsulAlertClient) getProfileForStatus(status string) string {
	// Appends s to folder.
	return c.getProfileForEntity("statu", status)
}

// GetProfileInfo returns profile info for check
func (c *ConsulAlertClient) GetProfileInfo(node, serviceID, checkID, status string) ProfileInfo {
	log.Println("Getting profile for node: ", node, " service: ", serviceID, " check: ", checkID)

	var profile string

	profile = c.getProfileForService(serviceID)
	if profile == "" {
		profile = c.getProfileForCheck(checkID)
	}
	if profile == "" {
		profile = c.getProfileForNode(node)
	}
	if profile == "" {
		profile = c.getProfileForStatus(status)
	}
	if profile == "" {
		profile = "default"
	}

	var checkProfile ProfileInfo
	key := fmt.Sprintf("consul-alerts/config/notif-profiles/%s", profile)
	log.Println("profile key: ", key)
	kvPair, _, _ := c.api.KV().Get(key, nil)
	if kvPair == nil {
		log.Println("profile key not found.")
		return checkProfile
	}

	if err := json.Unmarshal(kvPair.Value, &checkProfile); err != nil {
		log.Error("Profile unmarshalling error: ", err.Error())
	} else {
		log.Println("Interval: ", checkProfile.Interval, " List: ", checkProfile.NotifList)
	}

	return checkProfile
}

// IsBlacklisted gets the blacklist status of check
func (c *ConsulAlertClient) IsBlacklisted(check *Check) bool {
	blacklistExist := func() bool {
		kvPairs, _, err := c.api.KV().List("consul-alerts/config/checks/blacklist/", nil)
		return len(kvPairs) != 0 && err == nil
	}

	node := check.Node
	nodeCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/nodes/%s", node)
	nodeBlacklisted := func() bool {
		return c.CheckKeyExists(nodeCheckKey) || c.CheckKeyMatchesRegexp("consul-alerts/config/checks/blacklist/nodes", node)
	}

	service := "_"
	serviceBlacklisted := func() bool { return false }
	if check.ServiceID != "" {
		service = check.ServiceID
		serviceCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/services/%s", service)

		serviceBlacklisted = func() bool {
			return c.CheckKeyExists(serviceCheckKey) || c.CheckKeyMatchesRegexp("consul-alerts/config/checks/blacklist/services", service)
		}
	}

	checkID := check.CheckID
	checkCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/checks/%s", checkID)

	checkBlacklisted := func() bool {
		return c.CheckKeyExists(checkCheckKey) || c.CheckKeyMatchesRegexp("consul-alerts/config/checks/blacklist/checks", checkID)
	}

	status := "_"
	statusBlacklisted := func() bool { return false }
	if check.Status != "" {
		status = check.Status
		statusCheckKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/status/%s", status)
		statusBlacklisted = func() bool {
			return c.CheckKeyExists(statusCheckKey) || c.CheckKeyMatchesRegexp("consul-alerts/config/checks/blacklist/status", status)
		}
	}

	singleKey := fmt.Sprintf("consul-alerts/config/checks/blacklist/single/%s/%s/%s", node, service, checkID)
	singleBlacklisted := func() bool { return c.CheckKeyExists(singleKey) }

	return blacklistExist() && (nodeBlacklisted() || serviceBlacklisted() || checkBlacklisted() || statusBlacklisted() || singleBlacklisted())
}

// GetChangeThreshold gets the node/service/check specific override for change threshold
func (c *ConsulAlertClient) GetChangeThreshold(check *Check) int {
	service := check.ServiceID
	if service == "" {
		service = "_"
	}
	// List from most specific to least specific
	keys := []string{
		fmt.Sprintf("consul-alerts/config/checks/single/%s/%s/%s/change-threshold", check.Node, service, check.CheckID),
		fmt.Sprintf("consul-alerts/config/checks/check/%s/change-threshold", check.CheckID),
		fmt.Sprintf("consul-alerts/config/checks/service/%s/change-threshold", service),
		fmt.Sprintf("consul-alerts/config/checks/node/%s/change-threshold", check.Node),
	}

	for _, key := range keys {
		log.Debugf("Checking key %s for change-threshold override", key)
		kvpair, _, err := c.api.KV().Get(key, nil)
		if kvpair == nil || err != nil {
			continue
		}
		if val, err := strconv.Atoi(string(kvpair.Value)); err == nil {
			log.Debugf("Found change-threshold override: %d", val)
			return val
		}
	}
	return -1
}

func (c *ConsulAlertClient) CheckKeyExists(key string) bool {
	kvpair, _, err := c.api.KV().Get(key, nil)
	return kvpair != nil && err == nil
}

func (c *ConsulAlertClient) CheckKeyMatchesRegexp(regexpKey string, key string) bool {
	kvPair, _, _ := c.api.KV().Get(regexpKey, nil)
	if kvPair != nil {
		var regexpList []string
		json.Unmarshal(kvPair.Value, &regexpList)
		for _, pattern := range regexpList {
			matched, err := regexp.MatchString(pattern, key)
			if err != nil {
				log.Printf("unable to match %s against pattern %s. Error: %s\n",
					key, pattern, err.Error())
			}
			if matched {
				return true
			}
		}
	}
	return false
}

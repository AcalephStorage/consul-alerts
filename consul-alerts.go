// Consul Alerts is a tool to send alerts when checks changes status.
// It is built on top of consul KV, Health, and watch features.
package main

import (
	"fmt"
	"os"
	"syscall"

	"net/http"
	"os/signal"

	"encoding/json"
	"io/ioutil"

	"github.com/AcalephStorage/consul-alerts/consul"
	"github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/docopt/docopt-go"
)

const version = "Consul Alerts 0.3.3"
const usage = `Consul Alerts.

Usage:
  consul-alerts start [--alert-addr=<addr>] [--consul-addr=<consuladdr>] [--consul-dc=<dc>] [--consul-acl-token=<token>] [--watch-checks] [--watch-events] [--log-level=<level>] [--config-file=<file>]
  consul-alerts watch (checks|event) [--alert-addr=<addr>] [--log-level=<level>]
  consul-alerts --help
  consul-alerts --version

Options:
  --consul-acl-token=<token>   The consul ACL token [default: ""].
  --alert-addr=<addr>          The address for the consul-alert api [default: localhost:9000].
  --consul-addr=<consuladdr>   The consul api address [default: localhost:8500].
  --consul-dc=<dc>             The consul datacenter [default: dc1].
  --log-level=<level>          Set the logging level - valid values are "debug", "info", "warn", and "err" [default: warn].
  --watch-checks               Run check watcher.
  --watch-events               Run event watcher.
  --help                       Show this screen.
  --version                    Show version.
  --config-file=<file>         Path to the configuration file in JSON format

`

type stopable interface {
	stop()
}

var consulClient consul.Consul

func main() {
	log.SetLevel(log.InfoLevel)
	args, _ := docopt.Parse(usage, nil, true, version, false)

	switch {
	case args["start"].(bool):
		daemonMode(args)
	case args["watch"].(bool):
		watchMode(args)
	}
}

func daemonMode(arguments map[string]interface{}) {

	// Define options before setting in either config file or on command line
	loglevelString := ""
	consulAclToken := ""
	consulAddr := ""
	consulDc := ""
	watchChecks := false
	watchEvents := false
	addr := ""
	var confData map[string]interface{}

	// This exists check only works for arguments with no default. arguments with defaults will always exist.
	// Because of this the current code overrides command line flags with config file options if set.
	if configFile, exists := arguments["--config-file"].(string); exists {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(file, &confData)
		if err != nil {
			log.Error(err)
		}
		log.Debug("Config data: ", confData)
	}

	if confData["log-level"] != nil {
		loglevelString = confData["log-level"].(string)
	} else {
		loglevelString = arguments["--log-level"].(string)
	}
	if confData["consul-acl-token"] != nil {
		consulAclToken = confData["consul-acl-token"].(string)
	} else {
		consulAclToken = arguments["--consul-acl-token"].(string)
	}
	if confData["consul-addr"] != nil {
		consulAddr = confData["consul-addr"].(string)
	} else {
		consulAddr = arguments["--consul-addr"].(string)
	}
	if confData["consul-dc"] != nil {
		consulDc = confData["consul-dc"].(string)
	} else {
		consulDc = arguments["--consul-dc"].(string)
	}
	if confData["alert-addr"] != nil {
		addr = confData["alert-addr"].(string)
	} else {
		addr = arguments["--alert-addr"].(string)
	}
	if confData["watch-checks"] != nil {
		watchChecks = confData["watch-checks"].(bool)
	} else {
		watchChecks = arguments["--watch-checks"].(bool)
	}
	if confData["watch-events"] != nil {
		watchEvents = confData["watch-events"].(bool)
	} else {
		watchEvents = arguments["--watch-events"].(bool)
	}

	if loglevelString != "" {
		loglevel, err := log.ParseLevel(loglevelString)
		if err == nil {
			log.SetLevel(loglevel)
		} else {
			log.Println("Log level not set:", err)
		}
	}

	url := fmt.Sprintf("http://%s/v1/info", addr)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 201 {
		version := resp.Header.Get("version")
		resp.Body.Close()
		log.Printf("consul-alert daemon already running version: %s", version)
		os.Exit(1)
	}

	consulClient, err = consul.NewClient(consulAddr, consulDc, consulAclToken)
	if err != nil {
		log.Println("Cluster has no leader or is unreacheable.", err)
		os.Exit(3)
	}

	hostname, _ := os.Hostname()

	log.Println("Consul ACL Token:", consulAclToken)
	log.Println("Consul Alerts daemon started")
	log.Println("Consul Alerts Host:", hostname)
	log.Println("Consul Agent:", consulAddr)
	log.Println("Consul Datacenter:", consulDc)

	leaderCandidate := startLeaderElection(consulAddr, consulDc, consulAclToken)
	notifEngine := startNotifEngine()

	ep := startEventProcessor()
	cp := startCheckProcessor(leaderCandidate, notifEngine)

	http.HandleFunc("/v1/info", infoHandler)
	http.HandleFunc("/v1/process/events", ep.eventHandler)
	http.HandleFunc("/v1/process/checks", cp.checkHandler)
	http.HandleFunc("/v1/health/wildcard", healthWildcardHandler)
	http.HandleFunc("/v1/health", healthHandler)
	go startAPI(addr)

	log.Println("Started Consul-Alerts API")

	if watchChecks {
		go runWatcher(consulAddr, consulDc, addr, loglevelString, "checks")
	}
	if watchEvents {
		go runWatcher(consulAddr, consulDc, addr, loglevelString, "event")
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-ch
	cleanup(notifEngine, cp, ep, leaderCandidate)
}

func startAPI(addr string) {
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("Error starting Consul-Alerts API", err)
		os.Exit(1)
	}
}

func watchMode(arguments map[string]interface{}) {
	loglevelString, _ := arguments["--log-level"].(string)

	if loglevelString != "" {
		loglevel, err := log.ParseLevel(loglevelString)
		if err == nil {
			log.SetLevel(loglevel)
		} else {
			log.Println("Log level not set:", err)
		}
	}

	checkMode := arguments["checks"].(bool)
	eventMode := arguments["event"].(bool)
	addr := arguments["--alert-addr"].(string)

	var watchType string
	switch {
	case checkMode:
		watchType = "checks"
	case eventMode:
		watchType = "events"
	}

	url := fmt.Sprintf("http://%s/v1/process/%s", addr, watchType)
	resp, err := http.Post(url, "text/json", os.Stdin)
	if err != nil {
		log.Println("consul-alert daemon is not running.", err)
		os.Exit(2)
	} else {
		resp.Body.Close()
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("version", version)
	w.WriteHeader(201)
}

func cleanup(stopables ...stopable) {
	log.Println("Shutting down...")
	for _, s := range stopables {
		s.stop()
	}
}

func builtinNotifiers() []notifier.Notifier {

	emailConfig := consulClient.EmailConfig()
	logConfig := consulClient.LogConfig()
	influxdbConfig := consulClient.InfluxdbConfig()
	slackConfig := consulClient.SlackConfig()
	pagerdutyConfig := consulClient.PagerDutyConfig()
	hipchatConfig := consulClient.HipChatConfig()
	opsgenieConfig := consulClient.OpsGenieConfig()
	awssnsConfig := consulClient.AwsSnsConfig()
	victoropsConfig := consulClient.VictorOpsConfig()

	notifiers := []notifier.Notifier{}
	if emailConfig.Enabled {
		emailNotifier := &notifier.EmailNotifier{
			Url:         emailConfig.Url,
			Port:        emailConfig.Port,
			Username:    emailConfig.Username,
			Password:    emailConfig.Password,
			SenderAlias: emailConfig.SenderAlias,
			SenderEmail: emailConfig.SenderEmail,
			Receivers:   emailConfig.Receivers,
			Template:    emailConfig.Template,
			ClusterName: emailConfig.ClusterName,
			NotifName:   "email",
			OnePerAlert: emailConfig.OnePerAlert,
			OnePerNode:  emailConfig.OnePerNode,
		}
		notifiers = append(notifiers, emailNotifier)
	}
	if logConfig.Enabled {
		logNotifier := &notifier.LogNotifier{
			LogFile:   logConfig.Path,
			NotifName: "log",
		}
		notifiers = append(notifiers, logNotifier)
	}
	if influxdbConfig.Enabled {
		influxdbNotifier := &notifier.InfluxdbNotifier{
			Host:       influxdbConfig.Host,
			Username:   influxdbConfig.Username,
			Password:   influxdbConfig.Password,
			Database:   influxdbConfig.Database,
			SeriesName: influxdbConfig.SeriesName,
			NotifName:  "influx",
		}
		notifiers = append(notifiers, influxdbNotifier)
	}
	if slackConfig.Enabled {
		slackNotifier := &notifier.SlackNotifier{
			ClusterName: slackConfig.ClusterName,
			Url:         slackConfig.Url,
			Channel:     slackConfig.Channel,
			Username:    slackConfig.Username,
			IconUrl:     slackConfig.IconUrl,
			IconEmoji:   slackConfig.IconEmoji,
			Detailed:    slackConfig.Detailed,
			NotifName:   "slack",
		}
		notifiers = append(notifiers, slackNotifier)
	}
	if pagerdutyConfig.Enabled {
		pagerdutyNotifier := &notifier.PagerDutyNotifier{
			ServiceKey: pagerdutyConfig.ServiceKey,
			ClientName: pagerdutyConfig.ClientName,
			ClientUrl:  pagerdutyConfig.ClientUrl,
			NotifName:  "pagerduty",
		}
		notifiers = append(notifiers, pagerdutyNotifier)
	}
	if hipchatConfig.Enabled {
		hipchatNotifier := &notifier.HipChatNotifier{
			ClusterName: hipchatConfig.ClusterName,
			RoomId:      hipchatConfig.RoomId,
			AuthToken:   hipchatConfig.AuthToken,
			BaseURL:     hipchatConfig.BaseURL,
			From:        hipchatConfig.From,
			NotifName:   "hipchat",
		}
		notifiers = append(notifiers, hipchatNotifier)
	}
	if opsgenieConfig.Enabled {
		opsgenieNotifier := &notifier.OpsGenieNotifier{
			ClusterName: opsgenieConfig.ClusterName,
			ApiKey:      opsgenieConfig.ApiKey,
			NotifName:   "opsgenie",
		}
		notifiers = append(notifiers, opsgenieNotifier)
	}
	if awssnsConfig.Enabled {
		awssnsNotifier := &notifier.AwsSnsNotifier{
			Region:    awssnsConfig.Region,
			TopicArn:  awssnsConfig.TopicArn,
			NotifName: "awssns",
		}
		notifiers = append(notifiers, awssnsNotifier)
	}
	if victoropsConfig.Enabled {
		victoropsNotifier := &notifier.VictorOpsNotifier{
			APIKey:     victoropsConfig.APIKey,
			RoutingKey: victoropsConfig.RoutingKey,
			NotifName:  "victorops",
		}
		notifiers = append(notifiers, victoropsNotifier)
	}

	return notifiers
}

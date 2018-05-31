// Consul Alerts is a tool to send alerts when checks changes status.
// It is built on top of consul KV, Health, and watch features.
package main

import (
	"fmt"
	"strings"
	"os"
	"syscall"

	"net/http"
	"crypto/tls"
	"os/signal"

	"encoding/json"
	"io/ioutil"

	"github.com/AcalephStorage/consul-alerts/consul"
	"github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/docopt/docopt-go"
)

const version = "Consul Alerts 0.5.0"
const usage = `Consul Alerts.

Usage:
  consul-alerts start [--alert-addr=<addr>] [--consul-addr=<consuladdr>] [--consul-dc=<dc>] [--consul-acl-token=<token>] [--consul-ignore-cert=<bool>] [--watch-checks] [--watch-events] [--log-level=<level>] [--config-file=<file>]
  consul-alerts watch (checks|event) [--alert-addr=<addr>] [--log-level=<level>]
  consul-alerts --help
  consul-alerts --version

Options:
  --consul-acl-token=<token>   The consul ACL token [default: ""].
  --alert-addr=<addr>          The address for the consul-alert api [default: localhost:9000].
	--consul-addr=<consuladdr>   The consul api address [default: localhost:8500].
	--consul-ignore-cert=<bool>  Ignore the consul addr https certificate error
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
	scheme := "http"
	ignoreCert := false
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
	if confData["consul-ignore-cert"] != nil {
		ignoreCert = confData["consul-ignore-cert"].(bool)
	} else {
		ignoreCert = arguments["--consul-ignore-cert"].(bool)
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
	
	if strings.Contains(addr, "https") {
		scheme := "https"
	}

	tr := &http.Transport {
		if (scheme == "https" && ignoreCert)
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := fmt.Sprintf("%s://%s/v1/info", scheme, addr)
	resp, err := client.Get(url)
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
  
	url := fmt.Sprintf("%s://%s/v1/process/%s", scheme, addr, watchType)
	resp, err := client.Post(url, "text/json", os.Stdin)
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

func builtinNotifiers() map[string]notifier.Notifier {

	emailNotifier := consulClient.EmailNotifier()
	logNotifier := consulClient.LogNotifier()
	influxdbNotifier := consulClient.InfluxdbNotifier()
	slackNotifier := consulClient.SlackNotifier()
	mattermostNotifier := consulClient.MattermostNotifier()
	mattermostWebhookNotifier := consulClient.MattermostWebhookNotifier()
	pagerdutyNotifier := consulClient.PagerDutyNotifier()
	hipchatNotifier := consulClient.HipChatNotifier()
	opsgenieNotifier := consulClient.OpsGenieNotifier()
	awssnsNotifier := consulClient.AwsSnsNotifier()
	victoropsNotifier := consulClient.VictorOpsNotifier()

	notifiers := map[string]notifier.Notifier{}
	if emailNotifier.Enabled {
		notifiers[emailNotifier.NotifierName()] = emailNotifier
	}
	if logNotifier.Enabled {
		notifiers[logNotifier.NotifierName()] = logNotifier
	}
	if influxdbNotifier.Enabled {
		notifiers[influxdbNotifier.NotifierName()] = influxdbNotifier
	}
	if slackNotifier.Enabled {
		notifiers[slackNotifier.NotifierName()] = slackNotifier
	}
	if mattermostNotifier.Enabled {
		notifiers[mattermostNotifier.NotifierName()] = mattermostNotifier
	}
	if mattermostWebhookNotifier.Enabled {
		notifiers[mattermostWebhookNotifier.NotifierName()] = mattermostWebhookNotifier
	}
	if pagerdutyNotifier.Enabled {
		notifiers[pagerdutyNotifier.NotifierName()] = pagerdutyNotifier
	}
	if hipchatNotifier.Enabled {
		notifiers[hipchatNotifier.NotifierName()] = hipchatNotifier
	}
	if opsgenieNotifier.Enabled {
		notifiers[opsgenieNotifier.NotifierName()] = opsgenieNotifier
	}
	if awssnsNotifier.Enabled {
		notifiers[awssnsNotifier.NotifierName()] = awssnsNotifier
	}

	if victoropsNotifier.Enabled {
		notifiers[victoropsNotifier.NotifierName()] = victoropsNotifier
	}

	return notifiers
}

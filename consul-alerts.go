// Consul Alerts is a tool to send alerts when checks changes status.
// It is built on top of consul KV, Health, and watch features.
package main

import (
	"fmt"
	"os"
	"syscall"

	"net/http"
	"os/signal"

	"github.com/AcalephStorage/consul-alerts/consul"
	"github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/darkcrux/consul-skipper"
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/docopt/docopt-go"
)

const version = "Consul Alerts 0.3.0"
const usage = `Consul Alerts.

Usage:
  consul-alerts start [--alert-addr=<addr>] [--consul-addr=<consuladdr>] [--consul-dc=<dc>] [--watch-checks] [--watch-events]
  consul-alerts watch (checks|event) [--alert-addr=<addr>]
  consul-alerts --help
  consul-alerts --version

Options:
  --alert-addr=<addr>          The address for the consul-alert api [default: localhost:9000].
  --consul-addr=<consuladdr>   The consul api address [default: localhost:8500].
  --consul-dc=<dc>             The consul datacenter [default: dc1].
  --watch-checks               Run check watcher.
  --watch-events               Run event watcher.
  --help                       Show this screen.
  --version                    Show version.

`

var consulClient consul.Consul
var leaderCandidate *skipper.Candidate

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
	addr := arguments["--alert-addr"].(string)

	url := fmt.Sprintf("http://%s/v1/info", addr)
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {
		version := resp.Header.Get("version")
		resp.Body.Close()
		log.Printf("consul-alert daemon already running version: %s", version)
		os.Exit(1)
	}

	consulAddr := arguments["--consul-addr"].(string)
	consulDc := arguments["--consul-dc"].(string)
	watchChecks := arguments["--watch-checks"].(bool)
	watchEvents := arguments["--watch-events"].(bool)

	consulClient, err = consul.NewClient(consulAddr, consulDc)
	if err != nil {
		log.Println("Cluster has no leader or is unreacheable.", err)
		os.Exit(3)
	}

	hostname, _ := os.Hostname()

	log.Println("Consul Alerts daemon started")
	log.Println("Consul Alerts Host:", hostname)
	log.Println("Consul Agent:", consulAddr)
	log.Println("Consul Datacenter:", consulDc)

	leaderCandidate = &skipper.Candidate{
		ConsulAddress:    consulAddr,
		ConsulDatacenter: consulDc,
		LeadershipKey:    "consul-alerts/leader",
	}
	leaderCandidate.RunForElection()

	if watchChecks {
		go runWatcher(consulAddr, consulDc, "checks")
	}
	if watchEvents {
		go runWatcher(consulAddr, consulDc, "event")
	}

	go processEvents()
	go processChecks()

	http.HandleFunc("/v1/info", infoHandler)
	http.HandleFunc("/v1/process/events", eventHandler)
	http.HandleFunc("/v1/process/checks", checkHandler)
	http.HandleFunc("/v1/health", healthHandler)
	go http.ListenAndServe(addr, nil)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	cleanup()
}

func watchMode(arguments map[string]interface{}) {
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
		log.Println("consul-alert daemon is not running.")
		os.Exit(2)
	} else {
		resp.Body.Close()
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(201)
	w.Header().Add("version", version)
}

func cleanup() {
	log.Println("Shutting down...")
	leaderCandidate.Resign()
	close(checksChannel)
	close(eventsChannel)
}

func builtinNotifiers() []notifier.Notifier {

	emailConfig := consulClient.EmailConfig()
	logConfig := consulClient.LogConfig()
	influxdbConfig := consulClient.InfluxdbConfig()
	slackConfig := consulClient.SlackConfig()
	pagerdutyConfig := consulClient.PagerDutyConfig()

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
		}
		notifiers = append(notifiers, emailNotifier)
	}
	if logConfig.Enabled {
		logNotifier := &notifier.LogNotifier{
			LogFile: logConfig.Path,
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
		}
		notifiers = append(notifiers, slackNotifier)
	}
	if pagerdutyConfig.Enabled {
		pagerdutyNotifier := &notifier.PagerDutyNotifier{
			ServiceKey: pagerdutyConfig.ServiceKey,
			ClientName: pagerdutyConfig.ClientName,
			ClientUrl:  pagerdutyConfig.ClientUrl,
		}
		notifiers = append(notifiers, pagerdutyNotifier)
	}

	return notifiers
}

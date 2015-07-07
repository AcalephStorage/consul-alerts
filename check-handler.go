package main

import (
	"time"

	"net/http"

	"github.com/AcalephStorage/consul-alerts/consul"
	"github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

var checksChannel = make(chan []consul.Check, 1)
var firstCheckRun = true

func checkHandler(w http.ResponseWriter, r *http.Request) {
	consulClient.LoadConfig()
	if firstCheckRun {
		log.Println("Now watching for health changes.")
		firstCheckRun = false
		w.WriteHeader(200)
		return
	}

	if !consulClient.ChecksEnabled() {
		log.Println("Checks handling disabled. Checks ignored.")
		w.WriteHeader(200)
		return
	}

	if len(checksChannel) == 1 {
		<-checksChannel
	}

	var checks []consul.Check
	toWatchObject(r.Body, &checks)
	go startProcess(checks)
	w.WriteHeader(200)
}

func startProcess(checks []consul.Check) {
	checksChannel <- checks
}

func processChecks(notifEngine *NotifEngine) {
	for {
		<-checksChannel

		// if there's no leader, let's retry for at least 30 seconds in 5 second intervals.
		retryCount := 0
		for !hasLeader() {
			if retryCount >= 6 {
				continue
			}
			log.Println("There is current no consul-alerts leader... waiting for one.")
			time.Sleep(5 * time.Second)
			retryCount++
		}

		if !leaderCandidate.leader {
			log.Println("Currently not the leader. Ignoring checks.")
			continue
		}

		log.Println("Running health check.")
		changeThreshold := consulClient.CheckChangeThreshold()
		for elapsed := 0; elapsed < changeThreshold; elapsed += 10 {
			consulClient.UpdateCheckData()
			time.Sleep(10 * time.Second)
		}
		consulClient.UpdateCheckData()
		log.Println("Processing health checks for notification.")
		alerts := consulClient.NewAlerts()
		if len(alerts) > 0 {
			notify(notifEngine, alerts)
		}
	}
}

func notify(notifEngine *NotifEngine, alerts []consul.Check) {
	messages := make([]notifier.Message, len(alerts))
	for i, alert := range alerts {
		messages[i] = notifier.Message{
			Node:      alert.Node,
			ServiceId: alert.ServiceID,
			Service:   alert.ServiceName,
			CheckId:   alert.CheckID,
			Check:     alert.Name,
			Status:    alert.Status,
			Output:    alert.Output,
			Notes:     alert.Notes,
			Timestamp: time.Now(),
		}
	}

	if len(messages) == 0 {
		log.Println("Nothing to notify.")
		return
	}

	notifEngine.queueMessages(messages)
}

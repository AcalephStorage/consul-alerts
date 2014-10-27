package main

import (
	"bytes"
	"log"
	"time"

	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/AcalephStorage/consul-alerts/consul"
	"github.com/AcalephStorage/consul-alerts/notifier"
)

var checksChannel chan []consul.Check = make(chan []consul.Check, 1)
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
		log.Println("replacing health check request queue.")
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

func processChecks() {
	for {
		<-checksChannel
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
			notify(alerts)
		}
	}
}

func notify(alerts []consul.Check) {
	messages := make([]notifier.Message, len(alerts))
	for i, alert := range alerts {
		messages[i] = notifier.Message{
			Node:      alert.Node,
			Service:   alert.ServiceName,
			Check:     alert.Name,
			Status:    alert.Status,
			Output:    alert.Output,
			Notes:     alert.Notes,
			Timestamp: time.Now(),
		}
	}
	for _, n := range builtinNotifiers() {
		n.Notify(messages)
	}
	for _, n := range consulClient.CustomNotifiers() {
		executeHealthNotifier(messages, n)
	}
}

func executeHealthNotifier(messages []notifier.Message, notifCmd string) {
	data, err := json.Marshal(&messages)
	if err != nil {
		log.Println("Unable to read messages: ", err)
		return
	}

	input := bytes.NewReader(data)
	output := new(bytes.Buffer)
	cmd := exec.Command(notifCmd)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Run(); err != nil {
		log.Println("error running notifier: ", err)
	} else {
		log.Println(">>> notification sent to:", notifCmd)
	}
	log.Println(output)

}

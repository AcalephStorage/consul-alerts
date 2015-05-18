package main

import (
	"bytes"

	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/AcalephStorage/consul-alerts/consul"

	log "github.com/Sirupsen/logrus"
)

var eventsChannel = make(chan []consul.Event)

var firstEventRun bool = true

func eventHandler(w http.ResponseWriter, r *http.Request) {
	consulClient.LoadConfig()
	if firstEventRun {
		log.Println("Now watching for events.")
		firstEventRun = false
		// set status to OK
		return
	}

	if !consulClient.EventsEnabled() {
		log.Println("Event handling disabled. Event ignored.")
		// set to OK?
		return
	}

	var events []consul.Event
	toWatchObject(r.Body, &events)
	eventsChannel <- events
	// set status to OK
}

func processEvents() {
	for {
		events := <-eventsChannel
		for _, event := range events {
			processEvent(event)
		}
	}
}

func processEvent(event consul.Event) {
	log.Println("----------------------------------------")
	log.Printf("Processing event %s:\n", event.ID)
	log.Println("----------------------------------------")
	eventHandlers := consulClient.EventHandlers(event.Name)
	for _, eventHandler := range eventHandlers {
		executeEventHandler(event, eventHandler)
	}
	log.Printf("Event Processed.\n\n")
}

func executeEventHandler(event consul.Event, eventHandler string) {

	data, err := json.Marshal(&event)
	if err != nil {
		log.Println("Unable to read event: ", event)
		// then what?
	}

	input := bytes.NewReader(data)
	output := new(bytes.Buffer)
	cmd := exec.Command(eventHandler)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Run(); err != nil {
		log.Println("error running handler: ", err)
	} else {
		log.Printf(">>> \n%s -> %s:\n %s\n", event.ID, eventHandler, output)
	}
}

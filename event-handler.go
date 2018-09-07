package main

import (
	"bytes"

	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/AcalephStorage/consul-alerts/consul"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type EventProcessor struct {
	inChan    chan []consul.Event
	closeChan chan struct{}
	firstRun  bool
}

func (ep *EventProcessor) start() {
	cleanup := false
	for !cleanup {
		select {
		case events := <-ep.inChan:
			ep.handleEvents(events)
		case <-ep.closeChan:
			cleanup = true
		}
	}
}

func (ep *EventProcessor) stop() {
	close(ep.closeChan)
}

func (ep *EventProcessor) handleEvents(events []consul.Event) {
	for _, event := range events {
		log.Println("----------------------------------------")
		log.Printf("Processing event %s:\n", event.ID)
		log.Println("----------------------------------------")
		eventHandlers := consulClient.EventHandlers(event.Name)
		for _, eventHandler := range eventHandlers {
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
		log.Printf("Event Processed.\n\n")
	}
}

func (ep *EventProcessor) eventHandler(w http.ResponseWriter, r *http.Request) {
	consulClient.LoadConfig()
	if ep.firstRun {
		log.Println("Now watching for events.")
		ep.firstRun = false
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
	ep.inChan <- events
	// set status to OK
}

func startEventProcessor() *EventProcessor {
	ep := &EventProcessor{
		inChan:    make(chan []consul.Event, 1),
		closeChan: make(chan struct{}),
		firstRun:  true,
	}
	go ep.start()
	return ep
}

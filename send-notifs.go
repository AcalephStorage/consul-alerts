package main

import (
	"bytes"
	"encoding/json"
	"github.com/imdario/mergo"
	"os/exec"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/AcalephStorage/consul-alerts/notifier"
)

// NotifEngine handles notifications.
//
// To start NotifEngine:
//     notifEngine := startNotifEngine()
//
// Tp stop NotifEngine (for cleanup):
//     notifEngine.stop()
//
type NotifEngine struct {
	inChan    chan notifier.Messages
	closeChan chan struct{}
}

func (n *NotifEngine) start() {
	cleanup := false
	for !cleanup {
		select {
		case messages := <-n.inChan:
			n.sendBuiltin(messages)
			n.sendCustom(messages)
		case <-n.closeChan:
			cleanup = true
		}
	}
}

func (n *NotifEngine) stop() {
	close(n.closeChan)
}

func (n *NotifEngine) queueMessages(messages notifier.Messages) {
	n.inChan <- messages
	log.Println("messages sent for notification")
}

func (n *NotifEngine) sendBuiltin(messages notifier.Messages) {
	log.Println("sendBuiltin running")

	defaultNotifiers := builtinNotifiers()
	messagesPerNotifier := make(map[notifier.Notifier]notifier.Messages)

	for _, m := range messages {
		// if notification list is empty -> notify by all the enabled notifiers
		if len(m.NotifList) == 0 {
			for _, notifier := range defaultNotifiers {
				messagesPerNotifier[notifier] = append(messagesPerNotifier[notifier], m)
			}
		}

		for notifName, enabled := range m.NotifList {
			// get the default notifier
			if defaultNotifier, defaultNotifierExists := defaultNotifiers[notifName]; defaultNotifierExists && enabled {
				notif := defaultNotifier
				if varOverride, varOverrideExists := m.VarOverrides.GetNotifier(notifName); varOverrideExists {
					err := mergo.MergeWithOverwrite(notif, varOverride)
					if err != nil {
						log.Error(err)
					}
				}
				messagesPerNotifier[notif] = append(messagesPerNotifier[notif], m)
			}
		}
	}

	for notifier, msgs := range messagesPerNotifier {
		notifier.Notify(msgs)
	}
}

func (n *NotifEngine) sendCustom(messages notifier.Messages) {
	for notifName, notifCmd := range consulClient.CustomNotifiers() {
		filteredMessages := make(notifier.Messages, 0)
		for _, m := range messages {
			if boolVal, exists := m.NotifList[notifName]; (exists && boolVal) || len(m.NotifList) == 0 {
				filteredMessages = append(filteredMessages, m)
			}
		}
		if len(filteredMessages) == 0 {
			continue
		}
		data, err := json.Marshal(&filteredMessages)
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
}

func startNotifEngine() *NotifEngine {
	notifEngine := &NotifEngine{
		inChan:    make(chan notifier.Messages),
		closeChan: make(chan struct{}),
	}
	go notifEngine.start()
	return notifEngine
}

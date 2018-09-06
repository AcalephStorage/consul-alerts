package main

import (
	"math"
	"time"

	"net/http"

	"github.com/cgetzen/consul-alerts/consul"
	"github.com/AcalephStorage/consul-alerts/notifier"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type CheckProcessor struct {
	inChan         chan []consul.Check
	closeChan      chan struct{}
	firstRun       bool
	notifEngine    *NotifEngine
	leaderElection *LeaderElection
}

func (c *CheckProcessor) start() {
	cleanup := false
	for !cleanup {
		select {
		case checks := <-c.inChan:
			c.handleChecks(checks)
		case <-c.closeChan:
			cleanup = true
		}
	}
}

func (c *CheckProcessor) stop() {
	close(c.closeChan)
}

func (c *CheckProcessor) reminderStart() {
	cleanup := false
	remindTicker := time.NewTicker(time.Second * 300).C
	for !cleanup {
		select {
		case <-remindTicker:
			c.reminderRun()
		case <-c.closeChan:
			cleanup = true
		}
	}
}

func (c *CheckProcessor) reminderRun() {
	if !c.leaderElection.leader {
		log.Println("Currently not the leader. Ignoring reminders.")
		return
	}
	log.Println("Running reminder check.")
	messages := consulClient.GetReminders()
	filteredMessages := make(notifier.Messages, 0)
	for _, message := range messages {
		check := &consul.Check{
			Node:        message.Node,
			CheckID:     message.CheckId,
			Name:        message.Check,
			Status:      message.Status,
			Notes:       message.Notes,
			Output:      message.Output,
			ServiceID:   message.ServiceId,
			ServiceName: message.Service,
		}
		if consulClient.IsBlacklisted(check) {
			log.Printf("%s:%s:%s is blacklisted, deleting reminder", check.Node, check.ServiceID, check.CheckID)
			consulClient.DeleteReminder(check.Node, check.CheckID)
			continue
		}
		duration := time.Since(message.RmdCheck)
		durMins := int(math.Ceil(duration.Minutes()))
		log.Println("Reminder message duration minutes: ", durMins)
		if durMins >= message.Interval {
			message.RmdCheck = time.Now()
			consulClient.SetReminder(message)
			filteredMessages = append(filteredMessages, message)
		}
	}
	if len(filteredMessages) > 0 {
		c.notifEngine.queueMessages(filteredMessages)
	}
}

func (c *CheckProcessor) handleChecks(checks []consul.Check) {
	consulClient.LoadConfig()

	retryCount := 0
	for !hasLeader() {
		if retryCount >= 6 {
			return
		}
		log.Println("There is current no consul-alerts leader... waiting for one.")
		time.Sleep(5 * time.Second)
		retryCount++
	}

	if !c.leaderElection.leader {
		log.Println("Currently not the leader. Ignoring checks.")
		return
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
		c.notify(alerts)
	}

}

func (c *CheckProcessor) notify(alerts []consul.Check) {
	messages := make([]notifier.Message, len(alerts))
	for i, alert := range alerts {
		profileInfo := consulClient.GetProfileInfo(alert.Node, alert.ServiceID, alert.CheckID, "none")
		messages[i] = notifier.Message{
			Node:         alert.Node,
			ServiceId:    alert.ServiceID,
			Service:      alert.ServiceName,
			CheckId:      alert.CheckID,
			Check:        alert.Name,
			Status:       alert.Status,
			Output:       alert.Output,
			Notes:        alert.Notes,
			Interval:     profileInfo.Interval,
			RmdCheck:     time.Now(),
			NotifList:    profileInfo.NotifList,
			VarOverrides: profileInfo.VarOverrides,
			Timestamp:    time.Now(),
		}
		if profileInfo.Interval > 0 {
			switch alert.Status {
			case "passing":
				consulClient.DeleteReminder(alert.Node, alert.CheckID)
			case "warning", "critical":
				consulClient.SetReminder(messages[i])
			}
		}
	}

	if len(messages) == 0 {
		log.Println("Nothing to notify.")
		return
	}

	c.notifEngine.queueMessages(messages)
}

func startCheckProcessor(leaderCandidate *LeaderElection, notifEngine *NotifEngine) *CheckProcessor {
	cp := &CheckProcessor{
		inChan:         make(chan []consul.Check, 1),
		closeChan:      make(chan struct{}),
		firstRun:       true,
		notifEngine:    notifEngine,
		leaderElection: leaderCandidate,
	}
	go cp.start()
	go cp.reminderStart()
	return cp
}

func (c *CheckProcessor) checkHandler(w http.ResponseWriter, r *http.Request) {
	consulClient.LoadConfig()
	if c.firstRun {
		log.Println("Now watching for health changes.")
		c.firstRun = false
		w.WriteHeader(200)
		return
	}

	if !consulClient.ChecksEnabled() {
		log.Println("Checks handling disabled. Checks ignored.")
		w.WriteHeader(200)
		return
	}

	if len(c.inChan) == 1 {
		<-c.inChan
	}

	var checks []consul.Check
	toWatchObject(r.Body, &checks)
	c.inChan <- checks
	w.WriteHeader(200)
}

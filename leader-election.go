package main

import (
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"os"
	"syscall"
	"time"

	"os/signal"

	consulapi "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/hashicorp/consul/api"
)

type LeaderElection struct {
	lock           *api.Lock
	cleanupChannel chan struct{}
	stopChannel    chan struct{}
	leader         bool
}

func (l *LeaderElection) start() {
	clean := false
	for !clean {
		select {
		case <-l.cleanupChannel:
			clean = true
		default:
			log.Debugln("Acquiring Leadership")
			intChan, _ := l.lock.Lock(l.stopChannel)
			if intChan != nil {
				log.Debugln("Leadership Acquired")
				l.leader = true
				<-intChan
				l.leader = false
				fmt.Println("Leadership Lost")
				l.lock.Unlock()
				l.lock.Destroy()
			}
		}
	}
}

func (l *LeaderElection) stop() {
	fmt.Println("Cleaning up leadership")
	l.cleanupChannel <- struct{}{}
	l.stopChannel <- struct{}{}
	l.lock.Unlock()
	l.lock.Destroy()
	l.leader = false
	fmt.Print("cleanup done")
}

func startLeaderElection(client *consulapi.Client) *LeaderElection {
	lock, _ := client.LockKey("consul-alerts/leader")

	leader := &LeaderElection{
		lock:           lock,
		cleanupChannel: make(chan struct{}, 1),
		stopChannel:    make(chan struct{}, 1),
	}

	go leader.start()

	return leader
}

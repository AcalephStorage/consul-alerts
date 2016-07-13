package notifier

import (
	"log"
	"os"
	"path"

	"github.com/vincentvu/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type LogNotifier struct {
	LogFile   string
	NotifName string
}

// NotifierName provides name for notifier selection
func (logNotifier *LogNotifier) NotifierName() string {
	return logNotifier.NotifName
}

//Notify sends messages to the endpoint notifier
func (logNotifier *LogNotifier) Notify(alerts Messages) bool {

	logrus.Println("logging messages...")

	logDir := path.Dir(logNotifier.LogFile)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logrus.Printf("unable to create directory for logfile: %v\n", err)
		return false
	}

	file, err := os.OpenFile(logNotifier.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logrus.Printf("unable to write to logfile: %v\n", err)
		return false
	}

	logger := log.New(file, "[consul-notifier] ", log.LstdFlags)
	for _, alert := range alerts {
		logger.Printf("Node=%s, Service=%s, Check=%s, Status=%s\n", alert.Node, alert.Service, alert.Check, alert.Status)
	}
	logrus.Println("Notifications logged.")
	return true
}

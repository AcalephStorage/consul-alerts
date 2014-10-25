package notifier

import (
	"log"
	"os"
	"path"
)

type LogNotifier struct {
	LogFile string
}

func (logNotifier *LogNotifier) Notify(alerts []Message) bool {

	log.Println("logging messages...")

	logDir := path.Dir(logNotifier.LogFile)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Printf("unable to create directory for logfile: %v\n", err)
		return false
	}

	file, err := os.OpenFile(logNotifier.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("unable to write to logfile: %v\n", err)
		return false
	}

	logger := log.New(file, "[consul-notifier] ", log.LstdFlags)
	for _, alert := range alerts {
		logger.Printf("Node=%s, Service=%s, Check=%s, Status=%s\n", alert.Node, alert.Service, alert.Check, alert.Status)
	}
	log.Println("Notifications logged.")
	return true
}

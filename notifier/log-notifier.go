package notifier

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

type LogNotifier struct {
	Enabled bool
	Path    string `json:"path"`
}

// NotifierName provides name for notifier selection
func (logNotifier *LogNotifier) NotifierName() string {
	return "log"
}

func (logNotifier *LogNotifier) Copy() Notifier {
	notifier := *logNotifier
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (logNotifier *LogNotifier) Notify(alerts Messages) bool {

	logrus.Println("logging messages...")

	logDir := path.Dir(logNotifier.Path)
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logrus.Error(fmt.Sprintf("unable to create directory for logfile: %v", err))
		return false
	}

	file, err := os.OpenFile(logNotifier.Path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logrus.Error(fmt.Sprintf("unable to write to logfile: %v", err))
		return false
	}

	logger := log.New(file, "[consul-notifier] ", log.LstdFlags)
	for _, alert := range alerts {
		logger.Printf("Node=%s, Service=%s, Check=%s, Status=%s\n", alert.Node, alert.Service, alert.Check, alert.Status)
	}
	logrus.Println("Notifications logged.")
	return true
}

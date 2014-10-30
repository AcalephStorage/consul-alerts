package notifier

import "time"

const (
	SYSTEM_HEALTHY  string = "HEALTHY"
	SYSTEM_UNSTABLE string = "UNSTABLE"
	SYSTEM_CRITICAL string = "CRITICAL"
)

type Message struct {
	Node      string
	Service   string
	Check     string
	Status    string
	Output    string
	Notes     string
	Timestamp time.Time
}

type Messages []Message

type Notifier interface {
	Notify(alerts Messages) bool
}

func (m Message) IsCritical() bool {
	return m.Status == "critical"
}

func (m Message) IsWarning() bool {
	return m.Status == "warning"
}

func (m Message) IsPassing() bool {
	return m.Status == "passing"
}

func (m Messages) Summary() (overallStatus string, pass, warn, fail int) {
	hasCritical := false
	hasWarnings := false
	for _, message := range m {
		switch {
		case message.IsCritical():
			hasCritical = true
			fail++
		case message.IsWarning():
			hasWarnings = true
			warn++
		case message.IsPassing():
			pass++
		}
	}
	if hasCritical {
		overallStatus = SYSTEM_CRITICAL
	} else if hasWarnings {
		overallStatus = SYSTEM_UNSTABLE
	} else {
		overallStatus = SYSTEM_HEALTHY
	}
	return
}

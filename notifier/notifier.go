package notifier

import "time"

type Message struct {
	Node      string
	Service   string
	Check     string
	Status    string
	Output    string
	Notes     string
	Timestamp time.Time
}

type Notifier interface {
	Notify(alerts []Message) bool
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

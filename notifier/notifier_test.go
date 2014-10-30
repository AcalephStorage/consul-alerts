package notifier

import "testing"

func TestMessageIsCritical(t *testing.T) {
	message := Message{Status: "critical"}
	if !message.IsCritical() {
		t.Error("message should be critical")
	}
}

func TestMessageIsWarning(t *testing.T) {
	message := Message{Status: "warning"}
	if !message.IsWarning() {
		t.Error("message should be warning")
	}
}

func TestMessageIsPassing(t *testing.T) {
	message := Message{Status: "passing"}
	if !message.IsPassing() {
		t.Error("message should be passing")
	}
}

func TestSystemIsHealthy(t *testing.T) {
	messages := Messages{
		Message{Status: "passing"},
		Message{Status: "passing"},
		Message{Status: "passing"},
		Message{Status: "passing"},
		Message{Status: "passing"},
	}
	stat, pass, warn, fail := messages.Summary()
	if stat != SYSTEM_HEALTHY || pass != 5 || warn != 0 || fail != 0 {
		t.Errorf("system should be healthy, status=%s, pass=%d, warn=%d, fail=%d", stat, pass, warn, fail)
	}
}

func TestSystemIsCritical(t *testing.T) {
	messages := Messages{
		Message{Status: "passing"},
		Message{Status: "passing"},
		Message{Status: "critical"},
		Message{Status: "passing"},
		Message{Status: "warning"},
	}
	stat, pass, warn, fail := messages.Summary()
	if stat != SYSTEM_CRITICAL || pass != 3 || warn != 1 || fail != 1 {
		t.Errorf("system should be critical, status=%s, pass=%d, warn=%d, fail=%d", stat, pass, warn, fail)
	}
}

func TestSystemIsUnstable(t *testing.T) {
	messages := Messages{
		Message{Status: "warning"},
		Message{Status: "passing"},
		Message{Status: "warning"},
		Message{Status: "warning"},
		Message{Status: "passing"},
	}
	stat, pass, warn, fail := messages.Summary()
	if stat != SYSTEM_UNSTABLE || pass != 2 || warn != 3 || fail != 0 {
		t.Errorf("system should be unstable, status=%s, pass=%d, warn=%d, fail=%d", stat, pass, warn, fail)
	}
}

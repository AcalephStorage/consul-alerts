package notifier

import (
	"os"
	"testing"
)

func TestNotifySNSWithDefaultTemplate(t *testing.T) {
	oldSendSNS := sendSNS

	defer func() {
		sendSNS = oldSendSNS
	}()

	expectedSubject := "CRITICAL--Fail: 1, Warn: 0, Pass: 0"
	expectedMessage := "\nsome node:some service:some check is critical."

	sendSNS = func(awssns *AwsSnsNotifier, subject string, message string) bool {
		if subject != expectedSubject {
			t.Errorf("expected subject to be %s, got %s", expectedSubject, subject)
		}
		if message != expectedMessage {
			t.Errorf("expected message to be %s, got %s", expectedMessage, message)
		}
		return true
	}

	notifier := AwsSnsNotifier{
		Enabled:  true,
		Region:   "some region",
		TopicArn: "some-arn",
	}

	messages := Messages{Message{
		Node:    "some node",
		Service: "some service",
		Check:   "some check",
		Status:  "critical",
	}}
	if !notifier.Notify(messages) {
		t.Error("Notify must return true")
	}
}

func TestNotifySNSWithCustomTemplate(t *testing.T) {
	oldSendSNS := sendSNS

	defer func() {
		sendSNS = oldSendSNS
	}()

	expectedSubject := "CRITICAL--Fail: 1, Warn: 0, Pass: 0"
	expectedMessage := "custom template: Failed: 1"

	sendSNS = func(awssns *AwsSnsNotifier, subject string, message string) bool {
		if subject != expectedSubject {
			t.Errorf("expected subject to be %s, got %s", expectedSubject, subject)
		}
		if message != expectedMessage {
			t.Errorf("expected message to be %s, got %s", expectedMessage, message)
		}
		return true
	}

	tmpfile, err := templateFile("custom template: Failed: {{ .FailCount }}")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	notifier := AwsSnsNotifier{
		Template: tmpfile.Name(),
		Enabled:  true,
		Region:   "some region",
		TopicArn: "some-arn",
	}

	messages := Messages{Message{
		Node:    "some node",
		Service: "some service",
		Check:   "some check",
		Status:  "critical",
	}}
	if !notifier.Notify(messages) {
		t.Error("Notify must return true")
	}
}

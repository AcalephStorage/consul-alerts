package notifier

import (
	"fmt"
	"net/smtp"
	"reflect"
	"strings"
	"testing"
)

func TestNotify(t *testing.T) {
	oldSendMail := sendMail
	defer func() {
		sendMail = oldSendMail
	}()

	host := "mailserver.localdomain"
	port := 123

	expectedAddr := fmt.Sprintf("%s:%d", host, port)
	expectedFrom := "sender@example.com"
	expectedTo := []string{"test1@example.com", "test2@example.com"}
	expectedMsg := `From: "Some Sender" <sender@example.com>
To: test1@example.com, test2@example.com
Subject: Some Cluster is HEALTHY
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";


<!DOCTYPE html>
`

	sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		if addr != expectedAddr {
			t.Errorf("expected %s, got %s", expectedAddr, addr)
		}

		if a == nil {
			t.Error("auth must not be null")
		}

		if from != expectedFrom {
			t.Errorf("expected %s, got %s", expectedFrom, from)
		}

		if !reflect.DeepEqual(to, expectedTo) {
			t.Errorf("expected %s, got %s", expectedTo, to)
		}

		stringMsg := string(msg)
		if !strings.HasPrefix(stringMsg, expectedMsg) {
			t.Errorf("expected message to start with\n\n%s\n\ngot\n\n%s", expectedMsg, stringMsg)
		}

		return nil
	}

	notifier := EmailNotifier{
		Username:    "some username",
		Password:    "some password",
		ClusterName: "Some Cluster",
		Url:         host,
		Port:        port,
		SenderEmail: expectedFrom,
		SenderAlias: "Some Sender",
		Receivers:   expectedTo,
	}

	if !notifier.Notify(make(Messages, 0)) {
		t.Error("Notify must return true")
	}
}

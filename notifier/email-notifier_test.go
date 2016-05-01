package notifier

import (
	"net/smtp"
	"testing"
)

func TestNotify_AuthMustBeNilIfNoUsernameIsProvided(t *testing.T) {
	auth := notifyAndReturnAuth(t, EmailNotifier{Password: "some password"})

	if auth != nil {
		t.Error("auth must be nil if username is nil")
	}
}

func TestNotify_AuthMustBeNilIfNoPasswordIsProvided(t *testing.T) {
	auth := notifyAndReturnAuth(t, EmailNotifier{Username: "some username"})

	if auth != nil {
		t.Error("auth must be nil if password is nil")
	}
}

func TestNotify_AuthMustNotBeNilIfUsernameAndPasswordAreProvided(t *testing.T) {
	auth := notifyAndReturnAuth(t, EmailNotifier{Username: "some username", Password: "some password"})

	if auth == nil {
		t.Error("auth must not be nil if both username and password are not nil")
	}
}

func notifyAndReturnAuth(t *testing.T, notifier EmailNotifier) smtp.Auth {
	oldSendMail := sendMail
	defer func() {
		sendMail = oldSendMail
	}()

	var passedAuth smtp.Auth

	sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		passedAuth = a
		return nil
	}

	alerts := make(Messages, 0)
	ret := notifier.Notify(alerts)

	if !ret {
		t.Error("Notify must return true")
	}

	return passedAuth
}

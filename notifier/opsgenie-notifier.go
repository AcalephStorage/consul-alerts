package notifier

import (
	"fmt"

	alerts "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/alerts"
	ogcli "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/client"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type OpsGenieNotifier struct {
	Enabled     bool
	ClusterName string `json:"cluster-name"`
	ApiKey      string `json:"api-key"`
}

// NotifierName provides name for notifier selection
func (opsgenie *OpsGenieNotifier) NotifierName() string {
	return "opsgenie"
}

func (opsgenie *OpsGenieNotifier) Copy() Notifier {
	notifier := *opsgenie
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (opsgenie *OpsGenieNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	client := new(ogcli.OpsGenieClient)
	client.SetApiKey(opsgenie.ApiKey)

	alertCli, cliErr := client.Alert()

	if cliErr != nil {
		log.Println("Opsgenie notification trouble with client")
		return false
	}

	ok := true
	for _, message := range messages {
		title := fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		alias := opsgenie.createAlias(message)
		content := fmt.Sprintf(header, opsgenie.ClusterName, overallStatus, fail, warn, pass)
		content += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		content += fmt.Sprintf("\n%s", message.Output)

		// create the alert
		switch {
		case message.IsCritical():
			ok = opsgenie.createAlert(alertCli, title, content, alias) && ok
		case message.IsWarning():
			ok = opsgenie.createAlert(alertCli, title, content, alias) && ok
		case message.IsPassing():
			ok = opsgenie.closeAlert(alertCli, alias) && ok
		default:
			ok = false
			log.Warn("Message was not either IsCritical, IsWarning or IsPasssing. No notification was sent for ", alias)
		}
	}
	return ok
}

func (opsgenie OpsGenieNotifier) createAlias(message Message) string {
	incidentKey := message.Node
	if message.ServiceId != "" {
		incidentKey += ":" + message.ServiceId
	}

	return incidentKey
}

func (opsgenie *OpsGenieNotifier) createAlert(alertCli *ogcli.OpsGenieAlertClient, message string, content string, alias string) bool {
	log.Debug(fmt.Sprintf("OpsGenieAlertClient.CreateAlert alias: %s", alias))

	req := alerts.CreateAlertRequest{
		Message:     message,
		Description: content,
		Alias:       alias,
		Source:      "consul",
		Entity:      opsgenie.ClusterName,
	}
	response, alertErr := alertCli.Create(req)

	if alertErr != nil {
		if response == nil {
			log.Warn("Opsgenie notification trouble. ", alertErr)
		} else {
			log.Warn("Opsgenie notification trouble. ", response.Status)
		}
		return false
	}

	log.Println("Opsgenie notification sent.")
	return true
}

func (opsgenie *OpsGenieNotifier) closeAlert(alertCli *ogcli.OpsGenieAlertClient, alias string) bool {
	log.Debug(fmt.Sprintf("OpsGenieAlertClient.CloseAlert alias: %s", alias))
	req := alerts.CloseAlertRequest{
		Alias:  alias,
		Source: "consul",
	}
	response, alertErr := alertCli.Close(req)

	if alertErr != nil {
		if response == nil {
			log.Warn("Opsgenie notification trouble. ", alertErr)
		} else {
			log.Warn("Opsgenie notification trouble. ", response.Status)
		}
		return false
	}

	log.Println("Opsgenie close alert sent.")
	return true
}

package notifier

import (
	"fmt"

	alerts "github.com/opsgenie/opsgenie-go-sdk/alerts"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"

	log "github.com/Sirupsen/logrus"
)

type OpsGenieNotifier struct {
	ClusterName string
	ApiKey   string
}

func (opsgenie *OpsGenieNotifier) Notify(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	client := new(ogcli.OpsGenieClient)
	client.SetApiKey(opsgenie.ApiKey)

	alertCli, cliErr := client.Alert()

	if cliErr != nil {
		log.Println("Opsgenie notification trouble with client")
		return false
	}

	for _, message := range messages {
		title := fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		content := fmt.Sprintf(header, opsgenie.ClusterName, overallStatus, fail, warn, pass)
		content += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
		content += fmt.Sprintf("\n%s", message.Output)

		// create the alert
		response, alertErr := opsgenie.Send(alertCli, title, content)

		if alertErr != nil {
			log.Println("Opsgenie notification trouble.", response.Status)
			return false
		}

	}

	log.Println("Opsgenie notification send.")
	return true
}

func (opsgenie *OpsGenieNotifier) Send(alertCli *ogcli.OpsGenieAlertClient, message string, content string) (*alerts.CreateAlertResponse, error) {
	req := alerts.CreateAlertRequest{
		Message:         message,
		Description:     content,
		Source:          "consul",
		Entity:          opsgenie.ClusterName,
	}
	return alertCli.Create(req)
}

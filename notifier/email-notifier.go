package notifier

import (
	"fmt"

	"net/smtp"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"strings"
)

var sendMail = smtp.SendMail

type EmailNotifier struct {
	ClusterName string `json:"cluster-name"`
	Enabled     bool
	Template    string   `json:"template"`
	Url         string   `json:"url"`
	Port        int      `json:"port"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	SenderAlias string   `json:"sender-alias"`
	SenderEmail string   `json:"sender-email"`
	Receivers   []string `json:"receivers"`
	OnePerAlert bool     `json:"one-per-alert"`
	OnePerNode  bool     `json:"one-per-node"`
}

// NotifierName provides name for notifier selection
func (emailNotifier *EmailNotifier) NotifierName() string {
	return "email"
}

func (emailNotifier *EmailNotifier) Copy() Notifier {
	notifier := *emailNotifier
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (emailNotifier *EmailNotifier) Notify(alerts Messages) bool {

	overAllStatus, pass, warn, fail := alerts.Summary()
	nodeMap := mapByNodes(alerts)

	var emailDataList []TemplateData

	if emailNotifier.OnePerAlert {
		log.Println("Going to send one email per alert")
		emailDataList = []TemplateData{}
		for _, check := range alerts {

			singleAlertChecks := make(Messages, 0)
			singleAlertChecks = append(singleAlertChecks, check)
			singleAlertMap := mapByNodes(singleAlertChecks)

			alertStatus, alertPassing, alertWarnings, alertFailures := singleAlertChecks.Summary()

			alertClusterName := emailNotifier.ClusterName + " " + check.Node + " - " + check.CheckId

			e := TemplateData{
				ClusterName:  alertClusterName,
				SystemStatus: alertStatus,
				FailCount:    alertFailures,
				WarnCount:    alertWarnings,
				PassCount:    alertPassing,
				Nodes:        singleAlertMap,
			}
			emailDataList = append(emailDataList, e)
		}
	} else if emailNotifier.OnePerNode {
		log.Println("Going to send one email per node")
		emailDataList = []TemplateData{}
		for nodeName, checks := range nodeMap {
			singleNodeMap := mapByNodes(checks)
			nodeStatus, nodePassing, nodeWarnings, nodeFailures := checks.Summary()

			nodeClusterName := emailNotifier.ClusterName + " " + nodeName

			e := TemplateData{
				ClusterName:  nodeClusterName,
				SystemStatus: nodeStatus,
				FailCount:    nodeFailures,
				WarnCount:    nodeWarnings,
				PassCount:    nodePassing,
				Nodes:        singleNodeMap,
			}
			emailDataList = append(emailDataList, e)
		}
	} else {
		log.Println("Going to send one email for many alerts")
		e := TemplateData{
			ClusterName:  emailNotifier.ClusterName,
			SystemStatus: overAllStatus,
			FailCount:    fail,
			WarnCount:    warn,
			PassCount:    pass,
			Nodes:        nodeMap,
		}

		emailDataList = []TemplateData{e}
	}

	success := true

	for _, e := range emailDataList {

		var renderedTemplate string
		var err error
		renderedTemplate, err = renderTemplate(e, emailNotifier.Template, defaultTemplate)

		if err != nil {
			log.Println("Template error, unable to send email notification: ", err)
			success = false
			continue
		}

		msg := fmt.Sprintf(`From: "%s" <%s>
To: %s
Subject: %s is %s
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

%s
`,
			emailNotifier.SenderAlias,
			emailNotifier.SenderEmail,
			strings.Join(emailNotifier.Receivers, ", "),
			e.ClusterName,
			e.SystemStatus,
			renderedTemplate)

		addr := fmt.Sprintf("%s:%d", emailNotifier.Url, emailNotifier.Port)
		auth := smtp.PlainAuth("", emailNotifier.Username, emailNotifier.Password, emailNotifier.Url)
		if err := sendMail(addr, auth, emailNotifier.SenderEmail, emailNotifier.Receivers, []byte(msg)); err != nil {
			log.Println("Unable to send notification:", err)
			continue
		}
		log.Println("Email notification sent.")
		success = success && true
	}

	return success
}

func mapByNodes(alerts Messages) map[string]Messages {
	nodeMap := make(map[string]Messages)
	for _, alert := range alerts {
		nodeName := alert.Node
		nodeChecks := nodeMap[nodeName]
		if nodeChecks == nil {
			nodeChecks = make(Messages, 0)
		}
		nodeChecks = append(nodeChecks, alert)
		nodeMap[nodeName] = nodeChecks
	}
	return nodeMap
}

var defaultTemplate string = `
<!DOCTYPE html>
<html lang="en">
	<head>
  		<title>{{ .ClusterName }}</title>
	</head>

	<body style="width:100% !important; min-width: 100%; -webkit-text-size-adjust:100%; -ms-text-size-adjust:100%; margin:0; padding:0; font-family: 'Helvetica', 'Arial', sans-serif; color: #000000;">

		<div style="margin-left: auto; margin-right: auto; width: 36em; padding: 10dp; font-weight: bold; color: #ffffff; background-color: {{ if .IsCritical }}#e13329{{ else if .IsWarning }}#eebb00{{ else if .IsPassing }}#24c75a{{ end }};">
			<div style="padding: 10px;">
				{{ .ClusterName }}
			</div>
		</div>

		<div style="margin-left: auto; margin-right: auto; width: 36em; margin-top: 10px; margin-bottom: 10px; padding: 10dp">
			<p>
			<span style="font-weight: bold; font-size: 1.05em;">System is {{ .SystemStatus }}</span>
			<br/>
			<span style="font-size: 0.9em;">The following nodes are currently experiencing issues:</span>
			<div style="font-size: 0.85em;">
				<div style="float: left; width: 33%;">
					<strong>Failed: </strong>
					<span>{{ .FailCount }}</span>
				</div>
				<div style="float: right; width: 33%;">
					<strong>Warning: </strong>
					<span>{{ .WarnCount }}</span>
				</div>
				<div style="display: inline-block; width: 33%;">
					<strong>Passed: </strong>
					<span>{{ .PassCount }}</span>
				</div>
			</div>
			</p>

		</div>

		{{ range $name, $checks := .Nodes }}
		<div style="margin-left: auto; margin-right: auto; width: 36em; padding-top: 5px; padding-bottom: 20px;">
			<div style="font-size: 1.1em;">
				<strong>Node: </strong>
				<strong>{{ $name }}</strong>
			</div>

			{{ range $check := $checks }}
			<div style="margin-top: 15px; padding: 10px; background-color: {{ if $check.IsCritical }}#e13329{{ else if $check.IsWarning }}#eebb00{{ else if $check.IsPassing }}#24c75a{{ end }};">
				<div style="font-weight: bold; font-size: 1.1em;">
					{{ with $check.Service }}
					{{ $check.Service }}:
					{{ end }}
					{{ $check.Check }}
				</div>
				<div style="font-size: 0.85em;">
					<strong>Since: </strong>
					<span>{{ $check.Timestamp }}</span>
				</div>
				{{ with $check.Notes }}
				<div style="padding-top: 15px;">
					<strong>Notes: </strong>
					<pre>{{ $check.Notes }}</pre>
				</div>
				{{end }}
				<div style="padding-top: 15px;">
					<strong>Output:</strong>
					<pre>{{ $check.Output }}</pre>
				</div>
			</div>
			{{ end }}

		</div>
		{{ end }}


	</body>

</html>
`

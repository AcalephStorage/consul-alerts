package notifier

import (
	"fmt"
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type AwsSnsNotifier struct {
	Enabled  bool
	Template string `json:"template"`
	Region   string `json:"region"`
	TopicArn string `json:"topic-arn"`
}

// NotifierName provides name for notifier selection
func (awssns *AwsSnsNotifier) NotifierName() string {
	return "awssns"
}

func (awssns *AwsSnsNotifier) Copy() Notifier {
	notifier := *awssns
	return &notifier
}

func (awssns *AwsSnsNotifier) Notify(messages Messages) bool {
	subject := awssns.makeSubject(messages)
	body := awssns.makeBody(messages)

	return sendSNS(awssns, subject, body)
}

func (awssns *AwsSnsNotifier) makeSubject(messages Messages) string {
	overallStatus, pass, warn, fail := messages.Summary()
	return fmt.Sprintf("%s--Fail: %d, Warn: %d, Pass: %d", overallStatus, fail, warn, pass)
}

func (awssns *AwsSnsNotifier) makeBody(messages Messages) string {
	overallStatus, pass, warn, fail := messages.Summary()
	t := TemplateData{
		ClusterName:  "no-cluster-name",
		SystemStatus: overallStatus,
		FailCount:    fail,
		WarnCount:    warn,
		PassCount:    pass,
		Nodes:        mapByNodes(messages),
	}

	body, err := renderTemplate(t, awssns.Template, snsDefaultTemplate)
	if err != nil {
		log.Println("Template error, unable to send email notification: ", err)
		return fmt.Sprintf("error rendering template %v", err)
	} else {
		return body
	}
}

var sendSNS = func(awssns *AwsSnsNotifier, subject string, message string) bool {
	svc := sns.New(session.New(&aws.Config{
		Region: aws.String(awssns.Region),
	}))

	params := &sns.PublishInput{
		Message: aws.String(message),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Key": {
				DataType:    aws.String("String"),
				StringValue: aws.String("String"),
			},
		},
		MessageStructure: aws.String("messageStructure"),
		Subject:          aws.String(subject),
		TopicArn:         aws.String(awssns.TopicArn),
	}

	resp, err := svc.Publish(params)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	log.Println(resp)

	return true
}

var snsDefaultTemplate string = `
{{ range $name, $checks := .Nodes }}{{ range $check := $checks }}{{ $name }}:{{$check.Service}}:{{$check.Check}} is {{$check.Status}}.{{ end }}{{ end }}`

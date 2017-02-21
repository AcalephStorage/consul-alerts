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
	subject := MakeSubject(messages)
	body := MakeBody(messages)

	return awssns.Send(subject, body)
}

func MakeSubject(messages Messages) string {
	overallStatus, pass, warn, fail := messages.Summary()
	return fmt.Sprintf("%s--Fail: %d, Warn: %d, Pass: %d", overallStatus, fail, warn, pass)
}

func MakeBody(messages Messages) string {
	body := ""
	for _, message := range messages {
		body += fmt.Sprintf("\n%s:%s:%s is %s.", message.Node, message.Service, message.Check, message.Status)
	}
	return body
}

func (awssns *AwsSnsNotifier) Send(subject string, message string) bool {
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

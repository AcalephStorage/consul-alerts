package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type SumologicNotifier struct {
	NotifName    string
	Enabled      bool
	CollectorUri string
}

type SumologicMessage struct {
	Status      string `json:"status"`
	NodeName    string `json:"node_name"`
	ServiceName string `json:"service_name,omitempty"`
	ServiceID   string `json:"service_id,omitempty"`
	CheckName   string `json:"check_name"`
	CheckID     string `json:"check_id"`
	Output      string `json:"output"`
	Notes       string `json:"notes"`
}

func (n *SumologicNotifier) NotifierName() string {
	return "sumologic"
}

func (n *SumologicNotifier) Copy() Notifier {
	notifier := *n
	return &notifier
}

func (n *SumologicNotifier) Notify(alerts Messages) bool {
	if n.CollectorUri == "" {
		log.Error("sumologic collector uri is not configured")
		return false
	}

	ok := true

	for _, message := range alerts {
		host := message.Node
		name := message.Node
		category := "node"
		if message.ServiceId != "" {
			name = message.Service
			category = "service"
		}

		event := SumologicMessage{
			Status:      message.Status,
			NodeName:    message.Node,
			ServiceName: message.Service,
			ServiceID:   message.ServiceId,
			CheckName:   message.Check,
			CheckID:     message.CheckId,
			Output:      message.Output,
			Notes:       message.Notes,
		}

		if err := n.sendMessage(host, category, name, event); err != nil {
			ok = false
			log.Error(err.Error())
			continue
		}
	}

	log.Println("Sumologic message sent.")
	return ok
}

/// Send well-constructed message to sumologic
///  https://help.sumologic.com/Send_Data/Sources/02Sources_for_Hosted_Collectors/HTTP_Source/Upload_Data_to_an_HTTP_Source
func (n *SumologicNotifier) sendMessage(host string, name string, category string, event SumologicMessage) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializing sumologic message to json: %s", err)
	}

	req, err := http.NewRequest("POST", n.CollectorUri, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error constructing sumologic request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Add("X-Sumo-Host", host)
	req.Header.Add("X-Sumo-Category", category)
	req.Header.Add("X-Sumo-Name", name)

	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("error posting sumologic message: %s", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("sumologic returned bad status code: %d", res.StatusCode)
	}
	return nil
}

// Package notifier manages notifications for consul-alerts
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// VictorOpsNotifier provides configuration options for the VictorOps notifier
type VictorOpsNotifier struct {
	Enabled    bool
	APIKey     string `json:"api-key"`
	RoutingKey string `json:"routing-key"`
}

// VictorOpsEvent represents the options we'll pass to the VictorOps API
type VictorOpsEvent struct {
	// Explicitly listed by http://victorops.force.com/knowledgebase/articles/Integration/Alert-Ingestion-API-Documentation/
	MessageType       string `json:"message_type"`
	EntityID          string `json:"entity_id"`
	Timestamp         uint32 `json:"timestamp"`
	StateMessage      string `json:"state_message"`
	MonitoringTool    string `json:"monitoring_tool"`
	EntityDisplayName string `json:"entity_display_name"`

	// Helpful fields from http://victorops.force.com/knowledgebase/articles/Getting_Started/Incident-Fields-Glossary/?l=en_US&fs=RelatedArticle
	HostName    string `json:"host_name"`
	MonitorName string `json:"monitor_name"`

	// VictorOps lets you add arbitrary fields to help custom notification logic, so we'll set
	// node, service, service ID, check, and check ID
	ConsulNode      string `json:"consul_node"`
	ConsulService   string `json:"consul_service,omitempty"`
	ConsulServiceID string `json:"consul_service_id,omitempty"`
	ConsulCheck     string `json:"consul_check"`
	ConsulCheckID   string `json:"consul_check_id"`
}

const monitoringToolName string = "consul"
const apiEndpointTemplate string = "https://alert.victorops.com/integrations/generic/20131114/alert/%s/%s"

// NotifierName provides name for notifier selection
func (vo *VictorOpsNotifier) NotifierName() string {
	return "victorops"
}

func (vo *VictorOpsNotifier) Copy() Notifier {
	notifier := *vo
	return &notifier
}

// Notify sends messages to the endpoint notifier
func (vo *VictorOpsNotifier) Notify(messages Messages) bool {
	endpoint := fmt.Sprintf(apiEndpointTemplate, vo.APIKey, vo.RoutingKey)

	ok := true

	for _, message := range messages {
		entityID := fmt.Sprintf("%s:", message.Node)
		entityDisplayName := entityID

		// This might be a node level check without an explicit service
		if message.ServiceId == "" {
			entityID += message.CheckId
			entityDisplayName += message.Check
		} else {
			entityID += message.ServiceId
			entityDisplayName += message.Service
		}

		var messageType string

		switch {
		case message.IsCritical():
			messageType = "CRITICAL"
		case message.IsWarning():
			messageType = "WARNING"
		case message.IsPassing():
			messageType = "RECOVERY"
		default:
			log.Warn(fmt.Sprintf("Message with status %s was neither critical, warning, nor passing, reporting to VictorOps as INFO", message.Status))
			messageType = "INFO"
		}

		// VictorOps automatically displays the entity display name in notifications and page SMSs / emails,
		// so for brevity we don't repeat it in the "StateMessage" field
		stateMessage := fmt.Sprintf("%s: %s\n%s", messageType, message.Notes, message.Output)

		event := VictorOpsEvent{
			MessageType:       messageType,
			EntityID:          entityID,
			Timestamp:         uint32(message.Timestamp.Unix()),
			StateMessage:      stateMessage,
			MonitoringTool:    monitoringToolName,
			EntityDisplayName: entityDisplayName,

			HostName:    message.Node,
			MonitorName: message.Check,

			ConsulNode:      message.Node,
			ConsulService:   message.Service,
			ConsulServiceID: message.ServiceId,
			ConsulCheck:     message.Check,
			ConsulCheckID:   message.CheckId,
		}

		eventJSON, jsonError := json.Marshal(event)

		if jsonError != nil {
			ok = false
			log.Error("Error JSON-ifying VictorOps alert. ", jsonError)
			continue
		}

		response, httpError := http.Post(endpoint, "application/json", bytes.NewBuffer(eventJSON))

		if httpError != nil {
			ok = false
			log.Error("Error hitting VictorOps API. ", httpError)
			continue
		}

		if response.StatusCode != 200 {
			ok = false
			log.Error(fmt.Sprintf("Expected VictorOps endpoint to return 200, but it returned %d", response.StatusCode))
			continue
		}
	}

	log.Println("VictorOps notification sent.")
	return ok
}

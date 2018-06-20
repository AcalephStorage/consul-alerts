package notifier

import (
	"reflect"
	"testing"
)

func TestILertEvents(t *testing.T) {
	ilert := &ILertNotifier{
		ApiKey:              "aPiKeY",
		IncidentKeyTemplate: "{{.Check}}:{{.Node}}:{{.Service}}",
	}

	messages := Messages{
		Message{Status: "passing", Node: "node1", Service: "service1", Check: "check1", Output: "OK"},
		Message{Status: "warning", Node: "node1", Service: "service1", Check: "check2", Output: "WARN"},
		Message{Status: "critical", Node: "node2", Service: "service2", Check: "check3", Output: "CRIT"},
	}

	expectedILertEvents := []iLertEvent{
		{
			ApiKey:      "aPiKeY",
			EventType:   "RESOLVE",
			Summary:     "check1:node1:service1 is now HEALTHY",
			Details:     "OK",
			IncidentKey: "check1:node1:service1",
		},
		{
			ApiKey:      "aPiKeY",
			EventType:   "ALERT",
			Summary:     "check3:node2:service2 is CRITICAL",
			Details:     "CRIT",
			IncidentKey: "check3:node2:service2",
		},
	}

	actualILertEvents := ilert.toILertEvents(messages)

	if !reflect.DeepEqual(expectedILertEvents, actualILertEvents) {
		t.Errorf("iLert event mapping failed, expected: %v, got: %v", expectedILertEvents, actualILertEvents)
	}
}

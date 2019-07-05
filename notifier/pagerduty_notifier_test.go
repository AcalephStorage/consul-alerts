package notifier

import (
	"fmt"
	"testing"
)

func TestPagerDuty_StandardConfig(t *testing.T) {
	emptyMessages := Messages{}

	pd := PagerDutyNotifier{
		ClientName: "pagerduty",
		ClientUrl:  "example.com",
		Enabled:    true,
		ServiceKey: "longdummykey",
	}

	success := pd.Notify(emptyMessages)

	if !success {
		t.Fatalf("Expected Notify success with regular config and zero messages")
	}
}

func TestPagerDuty_WithOptionalConfigs(t *testing.T) {
	emptyMessages := Messages{}

	suite := []struct {
		name              string
		maxRetry          int
		retryBaseInterval int
	}{
		{"with MaxRetry", 5, 0},
		{"with RetryBaseInterval", 0, 30},
		{"with both MaxRetry RetryBaseInterval", 5, 30},
	}

	for _, tc := range suite {
		t.Run(tc.name, func(t *testing.T) {
			pd := PagerDutyNotifier{
				ClientName: "pagerduty",
				ClientUrl:  "example.com",
				Enabled:    true,
				ServiceKey: "longdummykey",
			}
			pd.MaxRetry = tc.maxRetry
			pd.RetryBaseInterval = tc.retryBaseInterval
			fmt.Printf("pd: %+v\n", pd)
			success := pd.Notify(emptyMessages)

			if !success {
				t.Fatalf("Expected Notify success for config '%s' and zero messages", tc.name)
			}
		})
	}
}

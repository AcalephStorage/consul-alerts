package gopherduty

import (
	"testing"
	"time"
)

type sampleDetail struct {
	Data1 string
	Data2 []string
}

func TestBackOffDelay(t *testing.T) {
	pd := &PagerDuty{
		MaxRetry:          3,
		RetryBaseInterval: 1,
	}

	delays := []int{1, 2, 4}

	for i := 0; i < pd.MaxRetry; i++ {
		now := time.Now()
		pd.retries = i
		pd.delayRetry()
		actualDelay := int(time.Since(now).Seconds())
		expectedDelay := delays[i]
		if actualDelay != expectedDelay {
			t.Errorf("expected delay is %d, actual delay is %d", expectedDelay, actualDelay)
		}
	}
}

func TestRetryOnRequest(t *testing.T) {
	pd := &PagerDuty{
		MaxRetry:          3,
		RetryBaseInterval: 1,
	}

	expectedRuntime := 7
	now := time.Now()
	response := pd.Trigger("", "", "", "", nil)
	actual := int(time.Since(now).Seconds())
	if !response.HasErrors() {
		t.Error("This should have been an error")
	}
	if actual < expectedRuntime {
		t.Errorf("Expected runtime is %d, actual is %d", expectedRuntime, actual)
	}

}

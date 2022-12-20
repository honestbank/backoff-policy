package backoff_policy

import (
	"time"

	"github.com/honestbank/backoff-policy/policies"
)

func NewExponentialBackoffPolicy(duration time.Duration, maxCount int) BackoffPolicy {
	return NewBackoff(policies.GetExponentialPolicy(1.1, duration, maxCount))
}

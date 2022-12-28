package backoff_policy

import (
	"time"

	"github.com/honestbank/backoff-policy/policies"
)

func NewExponentialBackoffPolicy(exponent float64, duration time.Duration, maxCount int) BackoffPolicy {
	return NewBackoff(policies.GetExponentialPolicy(exponent, duration, maxCount))
}

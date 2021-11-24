package backoff_policy

import (
	"time"

	"github.com/honestbank/backoff-policy/policies"
)

type backoff struct {
	count  int
	policy policies.Policy
}

func (b *backoff) MarkSuccess() {
	if b.count >= 1 {
		b.count -= 1
	}
}

func (b *backoff) MarkFailure() {
	b.count += 1
}

func (b *backoff) Execute(cb func(marker Marker)) {
	duration := b.policy(b.count)
	time.Sleep(duration)
	cb(b)
}

func NewBackoff(policy policies.Policy) BackoffPolicy {
	return &backoff{
		count:  0,
		policy: policy,
	}
}

package backoff_policy

import (
	"context"
	"time"

	"github.com/honestbank/backoff-policy/policies"
)

var _ ContextBackoffPolicy = (*backoff)(nil)

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

func (b *backoff) ExecuteWithContext(ctx context.Context, cb func(marker Marker)) {
	if duration := b.policy(b.count); duration > 0 {
		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
		}
	}
	cb(b)
}

func NewBackoff(policy policies.Policy) BackoffPolicy {
	return &backoff{
		count:  0,
		policy: policy,
	}
}

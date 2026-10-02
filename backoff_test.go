package backoff_policy_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	backoff_policy "github.com/honestbank/backoff-policy"
)

func TestNewBackoffSleepsAccordingToGivenPolicy(t *testing.T) {
	p := backoff_policy.NewBackoff(func(count int) time.Duration {
		return time.Duration(int(time.Millisecond) * 10 * count)
	})
	start := time.Now()
	p.Execute(func(marker backoff_policy.Marker) {
		marker.MarkFailure()
		marker.MarkFailure()
	})
	firstPass := time.Now()
	p.Execute(func(marker backoff_policy.Marker) {
		marker.MarkSuccess()
	})
	secondPass := time.Now()
	assert.Equal(t, int64(0), firstPass.Sub(start).Milliseconds())
	assert.Equal(t, int64(20), secondPass.Sub(start).Milliseconds())
}

func TestExecuteWithContext(t *testing.T) {
	// waitAfter waits for wait after 1 failure or more, and does not wait with no failures.
	waitAfter := func(wait time.Duration) func(count int) time.Duration {
		return func(count int) time.Duration {
			if count == 0 {
				return 0
			}

			return wait
		}
	}
	fail := func(marker backoff_policy.Marker) { marker.MarkFailure() }

	t.Run("NewBackoff and NewExponentialBackoffPolicy implement ContextBackoffPolicy", func(t *testing.T) {
		assert.Implements(t, (*backoff_policy.ContextBackoffPolicy)(nil), backoff_policy.NewBackoff(waitAfter(0)))
		assert.Implements(t, (*backoff_policy.ContextBackoffPolicy)(nil), backoff_policy.NewExponentialBackoffPolicy(time.Millisecond, 1))
	})

	t.Run("waits as Execute does", func(t *testing.T) {
		p := backoff_policy.NewBackoff(waitAfter(50 * time.Millisecond)).(backoff_policy.ContextBackoffPolicy)
		p.ExecuteWithContext(context.Background(), fail)

		start := time.Now()
		p.ExecuteWithContext(context.Background(), func(marker backoff_policy.Marker) { marker.MarkSuccess() })
		assert.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond)

		start = time.Now()
		p.ExecuteWithContext(context.Background(), func(marker backoff_policy.Marker) {})
		assert.Less(t, time.Since(start), 50*time.Millisecond)
	})

	t.Run("a done context ends the wait, and cb still runs", func(t *testing.T) {
		p := backoff_policy.NewBackoff(waitAfter(time.Hour)).(backoff_policy.ContextBackoffPolicy)
		p.ExecuteWithContext(context.Background(), fail)
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		called := false
		start := time.Now()
		p.ExecuteWithContext(ctx, func(marker backoff_policy.Marker) { called = true })

		assert.True(t, called)
		assert.Less(t, time.Since(start), 5*time.Second)
	})
}

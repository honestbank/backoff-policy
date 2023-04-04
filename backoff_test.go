package backoff_policy_test

import (
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

	// ideally secondPass.Sub(start).Milliseconds() should return 20,
	// but looks like an extra millisecond has passed between multiple calls
	assert.Equal(t, int64(21), secondPass.Sub(start).Milliseconds())
}

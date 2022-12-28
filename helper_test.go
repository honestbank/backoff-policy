package backoff_policy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	backoff_policy "github.com/honestbank/backoff-policy"
)

func TestNewExponentialBackoffPolicy(t *testing.T) {
	p := backoff_policy.NewExponentialBackoffPolicy(1, time.Millisecond*100, 2)
	assert.Panics(t, func() {
		p.Execute(func(marker backoff_policy.Marker) {
			panic("test panic")
		})
	})
}

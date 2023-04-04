package policies_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/backoff-policy/policies"
)

func TestGetConstantPolicy(t *testing.T) {
	t.Run("returns correct duration based on count", func(t *testing.T) {
		policy := policies.GetConstantPolicy(time.Millisecond*100)
		assert.Equal(t, time.Millisecond*0, policy(0))
		assert.Equal(t, time.Millisecond*100, policy(1))
		assert.Equal(t, time.Millisecond*100, policy(2))
		assert.Equal(t, time.Millisecond*100, policy(3))
	})
	t.Run("value is capped at max count", func(t *testing.T) {
		policy := policies.GetConstantPolicy(time.Millisecond*100)
		assert.Equal(t, policy(10), policy(10))
		assert.Equal(t, policy(10), policy(11))
		assert.Equal(t, policy(10), policy(12))
		assert.Equal(t, policy(10), policy(13))
	})
}

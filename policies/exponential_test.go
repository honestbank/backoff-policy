package policies_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/backoff-policy/policies"
)

func TestGetExponentialPolicy(t *testing.T) {
	t.Run("returns correct duration based on count", func(t *testing.T) {
		policy := policies.GetExponentialPolicy(2.0, time.Millisecond*100, 10)
		assert.Equal(t, time.Millisecond*0, policy(0))
		assert.Equal(t, time.Millisecond*200, policy(1))
		assert.Equal(t, time.Millisecond*400, policy(2))
		assert.Equal(t, time.Millisecond*3200, policy(5))
		assert.Equal(t, time.Millisecond*25600, policy(8))
		assert.Equal(t, time.Millisecond*102400, policy(10))
	})
	t.Run("value is capped at max count", func(t *testing.T) {
		policy := policies.GetExponentialPolicy(2.0, time.Millisecond*100, 10)
		assert.Equal(t, policy(10), policy(10))
		assert.Equal(t, policy(10), policy(11))
		assert.Equal(t, policy(10), policy(12))
		assert.Equal(t, policy(10), policy(13))
	})
}

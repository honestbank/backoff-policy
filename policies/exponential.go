package policies

import (
	"math"
	"time"
)

func GetExponentialPolicy(exponent float64, duration time.Duration, max int) Policy {
	return func(count int) time.Duration {
		if count == 0 {
			return 0
		}
		if count > max {
			count = max
		}
		multiplier := math.Pow(exponent, float64(count))

		return time.Duration(int64(float64(duration) * multiplier))
	}
}

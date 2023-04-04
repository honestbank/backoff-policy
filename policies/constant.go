package policies

import "time"

func GetConstantPolicy(duration time.Duration, max int) Policy {
	return GetExponentialPolicy(1, duration, max)
}

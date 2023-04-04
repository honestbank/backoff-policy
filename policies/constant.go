package policies

import "time"

func GetConstantPolicy(duration time.Duration) Policy {
	return func(count int) time.Duration {
		if count == 0 {
			return 0
		}

		return duration
	}
}

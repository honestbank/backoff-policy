# backoff-policy
In order to help developers wait if things are failing, this library aims to provide some useful methods.

The idea is that, you start a policy, in a loop, you do some work, if this fails, you mark it failed, if it passes, you mark it success.

This library will take care of delaying the work.

Code ideation:

```go
package main

import (
	"time"
	"github.com/honestbank/backoff-policy"
)

func process() error {
	// do work?
	return nil
}

func main() {
	policy := backoff_policy.NewExponentialBackoffPolicy(1.1, time.Millisecond * 100, 5)
	queue := getQueue()
	for ; ; {
		policy.Execute(func(marker backoff_policy.Marker) {
			message := queue.GetNewMessage()
			err := process(message)
			if err != nil {
				marker.MarkFailure()
				// log something, obviously...
				return
			}
			// queue.Ack(message)
			marker.MarkSuccess()
		})
	}
}
```

package backoff_policy

import "context"

type Marker interface {
	MarkSuccess()
	MarkFailure()
}
type BackoffPolicy interface {
	Execute(cb func(marker Marker))
}

// ContextBackoffPolicy is a BackoffPolicy whose wait a done context can end. The policies of NewBackoff and
// NewExponentialBackoffPolicy implement it, so a caller can check for it with a type assertion.
type ContextBackoffPolicy interface {
	BackoffPolicy
	// ExecuteWithContext waits as Execute does, but a done ctx ends the wait. cb then runs at once, so a caller
	// that must not run after ctx is done checks ctx.Err() in cb.
	ExecuteWithContext(ctx context.Context, cb func(marker Marker))
}

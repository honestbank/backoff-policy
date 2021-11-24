package backoff_policy

type Marker interface {
	MarkSuccess()
	MarkFailure()
}
type BackoffPolicy interface {
	Execute(cb func(marker Marker))
}

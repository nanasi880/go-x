package context

import "context"

// IsDone returns ctx is already done or else.
func IsDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// IsError returns true, if err is context.Canceled or context.DeadlineExceeded.
func IsError(err error) bool {
	switch err {
	case context.Canceled, context.DeadlineExceeded:
		return true
	default:
		return false
	}
}

// NewCancel is shorthand of context.WithCancel(context.Background()).
func NewCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

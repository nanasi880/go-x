package time

import (
	"context"
	"time"
)

// Sleep is sleep current goroutine.
func Sleep(ctx context.Context, d time.Duration) error {

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

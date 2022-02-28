package timeutil

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

// FixedZone returns a Location that always uses the given zone name and offset (east of UTC).
func FixedZone(name string, offset time.Duration) *time.Location {
	return time.FixedZone(name, int(offset.Seconds()))
}

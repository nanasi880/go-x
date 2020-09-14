package time_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	time2 "go.nanasi880.dev/x/time"
)

func TestSleep(t *testing.T) {

	type contexts struct {
		ctx    context.Context
		cancel context.CancelFunc
	}

	newContexts := func(ctx context.Context, cancelFunc context.CancelFunc) contexts {
		return contexts{
			ctx:    ctx,
			cancel: cancelFunc,
		}
	}

	ctx := context.Background()
	testCases := []struct {
		c       contexts
		sleepD  time.Duration
		cancelD time.Duration
	}{
		{
			c:       newContexts(ctx, nil),
			sleepD:  time.Second,
			cancelD: 0,
		},
		{
			c:       newContexts(context.WithCancel(ctx)),
			sleepD:  time.Second,
			cancelD: 500 * time.Millisecond,
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			wg := new(sync.WaitGroup)

			if tc.cancelD > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					time.Sleep(tc.cancelD)
					tc.c.cancel()
				}()
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				begin := time.Now()
				err := time2.Sleep(tc.c.ctx, tc.sleepD)
				end := time.Now()

				if tc.cancelD > 0 && err == nil {
					t.Errorf("")
					return
				}
				if tc.cancelD == 0 && err != nil {
					t.Errorf("%v", err)
					return
				}

				actualD := tc.sleepD
				if tc.cancelD > 0 {
					actualD = tc.cancelD
				}

				if end.Sub(begin) < actualD {
					t.Errorf("want: %v got: %v", actualD, end.Sub(begin))
				}
			}()

			wg.Wait()
		})
	}
}

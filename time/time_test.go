package time_test

import (
	"sync"
	"testing"
	"time"
	_ "time/tzdata"

	"go.nanasi880.dev/x/context"
	xtesting "go.nanasi880.dev/x/internal/testing"
	xtime "go.nanasi880.dev/x/time"
)

func TestSleep(t *testing.T) {
	testCases := []struct {
		ctx            context.CancelableContext
		sleepDuration  time.Duration
		cancelDuration time.Duration
	}{
		{
			ctx:            context.NewCancelableContext(xtesting.Context(t)),
			sleepDuration:  time.Second,
			cancelDuration: 0,
		},
		{
			ctx:            context.NewCancelableContext(xtesting.Context(t)),
			sleepDuration:  time.Second,
			cancelDuration: 500 * time.Millisecond,
		},
	}
	defer func() {
		for _, suite := range testCases {
			suite.ctx.Cancel()
		}
	}()

	for suiteNo, suite := range testCases {
		suiteNo, suite := suiteNo, suite

		wg := new(sync.WaitGroup)
		go func() {
			wg.Add(1)
			if suite.cancelDuration > 0 {
				time.Sleep(suite.cancelDuration)
				suite.ctx.Cancel()
			}
		}()

		wg.Wait()

		begin := time.Now()
		err := xtime.Sleep(suite.ctx, suite.sleepDuration)
		duration := time.Since(begin)

		if suite.cancelDuration > 0 && err == nil {
			xtesting.Failf(t, "suiteNo:%d suite.cancelDuration > 0 && err == nil", suiteNo)
			continue
		}
		if suite.cancelDuration == 0 && err != nil {
			xtesting.Failf(t, "suiteNo:%d suite.cancelDuration == 0 && err != nil", suiteNo)
			continue
		}

		actualDuration := suite.sleepDuration
		if suite.cancelDuration > 0 {
			actualDuration = suite.cancelDuration
		}
		if duration < actualDuration {
			xtesting.Failf(t, "suiteNo:%d duration:%v actual:%v", suiteNo, duration, actualDuration)
		}
	}
}

func TestFixedZone(t *testing.T) {
	const (
		minute = 60
		hour   = minute * 60
	)
	zone1 := time.FixedZone("Asia/Tokyo", 9*hour)
	zone2 := xtime.FixedZone("Asia/Tokyo", 9*time.Hour)

	now := time.Now()
	t1 := now.In(zone1)
	t2 := now.In(zone2)

	eq := t1.String() == t2.String()
	if !eq {
		xtesting.Fail(t, t1, t2)
	}
}

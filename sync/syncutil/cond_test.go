package syncutil_test

import (
	"context"
	"sync"
	"testing"
	"time"

	xsync "go.nanasi880.dev/x/sync/syncutil"
)

func TestCond_Signal(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cond := xsync.NewCond(new(sync.Mutex))

	go func() {
		time.Sleep(100 * time.Millisecond)
		cond.Signal()
		cancel()
	}()

	cond.L.Lock()
	defer cond.L.Unlock()

	err := cond.Wait(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkCond(b *testing.B) {
	b.Run("sync.Cond.Signal()", func(b *testing.B) {
		b.ReportAllocs()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cond := sync.NewCond(new(sync.Mutex))

		go func() {
			time.Sleep(time.Second)
			b.ResetTimer()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cond.Signal()
				}
			}
		}()
		for i := 0; i < b.N; i++ {
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
		}
	})
	b.Run("xsync.Cond.Signal()", func(b *testing.B) {
		b.ReportAllocs()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cond := xsync.NewCond(new(sync.Mutex))

		go func() {
			time.Sleep(time.Second)
			b.ResetTimer()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cond.Signal()
				}
			}
		}()
		background := context.Background()
		for i := 0; i < b.N; i++ {
			cond.L.Lock()
			_ = cond.Wait(background)
			cond.L.Unlock()
		}
	})
	b.Run("sync.Cond.Broadcast()", func(b *testing.B) {
		b.ReportAllocs()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cond := sync.NewCond(new(sync.Mutex))

		go func() {
			time.Sleep(time.Second)
			b.ResetTimer()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cond.Broadcast()
				}
			}
		}()
		for i := 0; i < b.N; i++ {
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
		}
	})
	b.Run("xsync.Cond.Broadcast()", func(b *testing.B) {
		b.ReportAllocs()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cond := xsync.NewCond(new(sync.Mutex))

		go func() {
			time.Sleep(time.Second)
			b.ResetTimer()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					cond.Broadcast()
				}
			}
		}()
		background := context.Background()
		for i := 0; i < b.N; i++ {
			cond.L.Lock()
			_ = cond.Wait(background)
			cond.L.Unlock()
		}
	})
}

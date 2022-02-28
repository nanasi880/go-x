package syncutil

import (
	"context"
	"sync"
	"sync/atomic"
)

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *Mutex or *RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
//
// A Cond must not be copied after first use.
type Cond struct {
	noCopy

	// L is held while observing or changing the condition
	L sync.Locker

	mutex  sync.Mutex
	init   int32
	notify chan struct{}
}

// NewCond returns a new Cond with Locker l.
func NewCond(l sync.Locker) *Cond {
	return &Cond{L: l}
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
func (c *Cond) Wait(ctx context.Context) error {
	c.initOnce()

	c.mutex.Lock()
	notify := c.notify
	c.mutex.Unlock()

	c.L.Unlock()
	defer c.L.Lock()

	select {
	case <-notify:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	c.initOnce()

	c.mutex.Lock()
	notify := c.notify
	c.mutex.Unlock()

	select {
	case notify <- struct{}{}:
		break
	default:
		break
	}
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	c.initOnce()

	c.mutex.Lock()
	defer c.mutex.Unlock()
	close(c.notify)
	c.notify = make(chan struct{})
}

func (c *Cond) initOnce() {
	if atomic.LoadInt32(&c.init) != 0 {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.init != 0 {
		return
	}

	c.init = 1
	c.notify = make(chan struct{})
}

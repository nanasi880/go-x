package context

import (
	"context"
	"time"
)

type CancelableContext interface {
	context.Context
	Cancel()
	IsDone() bool
}

type cancelableContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *cancelableContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *cancelableContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *cancelableContext) Err() error {
	return c.ctx.Err()
}

func (c *cancelableContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *cancelableContext) Cancel() {
	c.cancel()
}

func (c *cancelableContext) IsDone() bool {
	return IsDone(c.ctx)
}

func NewCancelable() CancelableContext {
	return NewCancelableContext(context.Background())
}

func NewCancelableContext(parent context.Context) CancelableContext {
	ctx, cancel := context.WithCancel(parent)
	return &cancelableContext{
		ctx:    ctx,
		cancel: cancel,
	}
}

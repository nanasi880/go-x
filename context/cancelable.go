package context

import (
	"context"
	"time"
)

// CancelableContext is cancelable context.
type CancelableContext interface {
	context.Context
	Cancel()
	IsDone() bool
}

type cancelableContext struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// Deadline is implementation of context.Context.
func (c *cancelableContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

// Done is implementation of context.Context.
func (c *cancelableContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err is implementation of context.Context.
func (c *cancelableContext) Err() error {
	return c.ctx.Err()
}

// Value is implementation of context.Context.
func (c *cancelableContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// Cancel is cancel this context.
func (c *cancelableContext) Cancel() {
	c.cancel()
}

// IsDone is shorthand of IsDone(ctx).
func (c *cancelableContext) IsDone() bool {
	return IsDone(c.ctx)
}

// Deconstruct is deconstruction this context to context.Context and context.CancelFunc.
func (c *cancelableContext) Deconstruct() (context.Context, context.CancelFunc) {
	return c.ctx, c.cancel
}

// NewCancelable is shorthand of NewCancelableContext(context.Background()).
func NewCancelable() CancelableContext {
	return NewCancelableContext(context.Background())
}

// NewCancelableContext is creating new cancelable context with parent context.
func NewCancelableContext(parent context.Context) CancelableContext {
	ctx, cancel := context.WithCancel(parent)
	return &cancelableContext{
		ctx:    ctx,
		cancel: cancel,
	}
}

package testing

import (
	"context"
	"testing"
)

// Context returns context bound testing handler.
func Context(tb testing.TB, parent ...context.Context) context.Context {
	tb.Helper()

	if len(parent) > 1 {
		tb.Fatal("`parent` only zero or one context")
	}

	var ctx context.Context
	switch len(parent) {
	case 0:
		ctx = context.Background()
	case 1:
		ctx = parent[0]
	default:
		tb.Fatal("`parent` only zero or one context")
	}

	ctx, cancel := context.WithCancel(ctx)
	tb.Cleanup(cancel)
	return ctx
}

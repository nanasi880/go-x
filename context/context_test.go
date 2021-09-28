package context_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	xcontext "go.nanasi880.dev/x/context"
)

func ExampleIsDone() {

	ctx, cancel := xcontext.NewCancel()

	fmt.Println(xcontext.IsDone(ctx))
	cancel()
	fmt.Println(xcontext.IsDone(ctx))

	// Output:
	// false
	// true
}

func ExampleIsError() {

	fmt.Println(xcontext.IsError(nil))
	fmt.Println(xcontext.IsError(context.Canceled))
	fmt.Println(xcontext.IsError(context.DeadlineExceeded))
	fmt.Println(xcontext.IsError(fmt.Errorf("context canceled")))
	// Output:
	// false
	// true
	// true
	// false
}

func ExampleCancelableContext() {
	ctx := xcontext.NewCancelable()

	fmt.Println(ctx.IsDone())
	ctx.Cancel()
	fmt.Println(ctx.IsDone())

	// Output:
	// false
	// true
}

func TestIsDone(t *testing.T) {
	ctx, cancel := xcontext.NewCancel()
	defer cancel()

	done := xcontext.IsDone(ctx)
	if done {
		t.Fatal()
	}

	cancel()
	done = xcontext.IsDone(ctx)
	if !done {
		t.Fatal()
	}
}

func TestIsError(t *testing.T) {
	if xcontext.IsError(nil) {
		t.Fatal()
	}
	if xcontext.IsError(errors.New("")) {
		t.Fatal()
	}
	if !xcontext.IsError(context.DeadlineExceeded) {
		t.Fatal()
	}
	if !xcontext.IsError(context.Canceled) {
		t.Fatal()
	}
}

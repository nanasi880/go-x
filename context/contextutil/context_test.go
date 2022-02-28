package contextutil_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.nanasi880.dev/x/context/contextutil"
)

func ExampleIsDone() {

	ctx, cancel := contextutil.NewCancel()

	fmt.Println(contextutil.IsDone(ctx))
	cancel()
	fmt.Println(contextutil.IsDone(ctx))

	// Output:
	// false
	// true
}

func ExampleIsError() {

	fmt.Println(contextutil.IsError(nil))
	fmt.Println(contextutil.IsError(context.Canceled))
	fmt.Println(contextutil.IsError(context.DeadlineExceeded))
	fmt.Println(contextutil.IsError(fmt.Errorf("context canceled")))
	// Output:
	// false
	// true
	// true
	// false
}

func ExampleCancelableContext() {
	ctx := contextutil.NewCancelable()

	fmt.Println(ctx.IsDone())
	ctx.Cancel()
	fmt.Println(ctx.IsDone())

	// Output:
	// false
	// true
}

func TestIsDone(t *testing.T) {
	ctx, cancel := contextutil.NewCancel()
	defer cancel()

	done := contextutil.IsDone(ctx)
	if done {
		t.Fatal()
	}

	cancel()
	done = contextutil.IsDone(ctx)
	if !done {
		t.Fatal()
	}
}

func TestIsError(t *testing.T) {
	if contextutil.IsError(nil) {
		t.Fatal()
	}
	if contextutil.IsError(errors.New("")) {
		t.Fatal()
	}
	if !contextutil.IsError(context.DeadlineExceeded) {
		t.Fatal()
	}
	if !contextutil.IsError(context.Canceled) {
		t.Fatal()
	}
}

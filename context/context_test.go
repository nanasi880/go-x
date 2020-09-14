package context_test

import (
	"context"
	"fmt"

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

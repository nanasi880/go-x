package mathutil_test

import (
	"math"
	"testing"

	"go.nanasi880.dev/x/internal/testing/testutil"
	"go.nanasi880.dev/x/math/mathutil"
)

func TestAbs32(t *testing.T) {
	inf := math.Inf(1)
	{
		v := math.Inf(1)
		abs := mathutil.Abs32(float32(v))
		if abs != float32(inf) {
			testutil.Fail(t, abs)
		}
	}
	{
		v := math.Inf(-1)
		abs := mathutil.Abs32(float32(v))
		if abs != float32(inf) {
			testutil.Fail(t, abs)
		}
	}
}

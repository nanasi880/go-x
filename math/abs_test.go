package math_test

import (
	"math"
	"testing"

	xtesting "go.nanasi880.dev/x/internal/testing"
	xmath "go.nanasi880.dev/x/math"
)

func TestAbs32(t *testing.T) {
	inf := math.Inf(1)
	{
		v := math.Inf(1)
		abs := xmath.Abs32(float32(v))
		if abs != float32(inf) {
			xtesting.Fail(t, abs)
		}
	}
	{
		v := math.Inf(-1)
		abs := xmath.Abs32(float32(v))
		if abs != float32(inf) {
			xtesting.Fail(t, abs)
		}
	}
}

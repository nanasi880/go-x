package errorutil_test

import (
	"testing"

	"go.nanasi880.dev/x/errors/errorutil"
)

func TestNew(t *testing.T) {
	err := errorutil.New("Hello")
	_ = err
}

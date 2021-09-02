package bytes_test

import (
	"fmt"
	"testing"

	"go.nanasi880.dev/x/bytes"
	xtesting "go.nanasi880.dev/x/internal/testing"
)

func TestSize_Format(t *testing.T) {
	testSuites := []struct {
		size   bytes.Size
		s      string
		format string
	}{
		{
			size:   0,
			s:      "0B",
			format: "%v",
		},
		{
			size:   bytes.KB - 1,
			s:      "999B",
			format: "%v",
		},
		{
			size:   -(bytes.KB - 1),
			s:      "-999B",
			format: "%v",
		},
		{
			size:   -1 * bytes.KB,
			s:      "-1.00KB",
			format: "%v",
		},
	}

	for i, suite := range testSuites {
		s := fmt.Sprintf(suite.format, suite.size)
		if s != suite.s {
			xtesting.Failf(t, "suite:%d want:%s got:%s", i, suite.s, s)
		}
	}
}

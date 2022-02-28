package byteutil_test

import (
	"fmt"
	"testing"

	"go.nanasi880.dev/x/bytes/byteutil"
	"go.nanasi880.dev/x/internal/testing/testutil"
)

func TestSize_Format(t *testing.T) {
	testSuites := []struct {
		size   byteutil.Size
		s      string
		format string
	}{
		{
			size:   0,
			s:      "0B",
			format: "%v",
		},
		{
			size:   byteutil.KB - 1,
			s:      "999B",
			format: "%v",
		},
		{
			size:   -(byteutil.KB - 1),
			s:      "-999B",
			format: "%v",
		},
		{
			size:   -1 * byteutil.KB,
			s:      "-1.00KB",
			format: "%v",
		},
	}

	for i, suite := range testSuites {
		s := fmt.Sprintf(suite.format, suite.size)
		if s != suite.s {
			testutil.Failf(t, "suite:%d want:%s got:%s", i, suite.s, s)
		}
	}
}

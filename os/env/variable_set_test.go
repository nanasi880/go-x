package env_test

import (
	"os"
	"testing"

	"go.nanasi880.dev/x/os/env"
)

func TestVariableSet(t *testing.T) {

	_ = os.Setenv("ENV_TEST_STRING", "=this is string")
	_ = os.Setenv("ENV_TEST_INT", "100")
	_ = os.Setenv("ENV_TEST_BOOL", "1")

	set := env.NewVariableSet("test", env.ContinueOnError)

	sp := set.String("ENV_TEST_STRING", "this is default", "usage")
	ip := set.Int("ENV_TEST_INT", -100, "usage")
	bp := set.Bool("ENV_TEST_BOOL", false, "usage")

	if err := set.Parse(os.Environ()); err != nil {
		t.Fatal(err)
	}

	if *sp != "=this is string" {
		t.Fatal(*sp)
	}
	if *ip != 100 {
		t.Fatal(*ip)
	}
	if *bp != true {
		t.Fatal(*bp)
	}
}

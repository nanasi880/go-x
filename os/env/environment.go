package env

import (
	"os"
)

// Default is the default set of environment variable set, parsed from os.Environment().
var Default = NewVariableSet("", ExitOnError)

// Usage returns Default.Usage().
var Usage = func() string {
	return Default.Usage()
}

// String defines a string value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the environment variable.
func String(name string, def string, usage string) *string {
	return Default.String(name, def, usage)
}

// StringVar defines a string environment variable with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the environment variable.
func StringVar(p *string, name string, def string, usage string) {
	Default.StringVar(p, name, def, usage)
}

// Int defines a int value with specified name, default value, and usage string.
// The return value is the address of a int variable that stores the value of the environment variable.
func Int(name string, def int, usage string) *int {
	return Default.Int(name, def, usage)
}

// IntVar defines a int environment variable with specified name, default value, and usage string.
// The argument p points to a int variable in which to store the value of the environment variable.
func IntVar(p *int, name string, def int, usage string) {
	Default.IntVar(p, name, def, usage)
}

// Bool defines a bool value with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the environment variable.
func Bool(name string, def bool, usage string) *bool {
	return Default.Bool(name, def, usage)
}

// BoolVar defines a bool environment variable with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the environment variable.
func BoolVar(p *bool, name string, def bool, usage string) {
	Default.BoolVar(p, name, def, usage)
}

// Parse parses the environment variables from os.Environ(). Must be called
// after all variables are defined and before variable are accessed by the program.
func Parse() {
	_ = Default.Parse(os.Environ())
}

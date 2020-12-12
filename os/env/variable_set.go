package env

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type variable struct {
	name  string
	value Value
	def   interface{}
	usage string
}

// A VariableSet represents a set of defined environment variable. The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
//
// Flag names must be unique within a FlagSet. An attempt to define a flag whose
// name is already in use will cause a panic.
type VariableSet struct {
	Usage       func() string
	name        string
	variables   map[string]*variable
	errHandling ErrorHandling
	output      io.Writer
}

// NewVariableSet returns a new, empty variable set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewVariableSet(name string, handling ErrorHandling) *VariableSet {
	set := &VariableSet{
		Usage:       nil,
		name:        name,
		errHandling: handling,
	}
	set.Usage = set.defaultUsage

	return set
}

// String defines a string value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the environment variable.
func (set *VariableSet) String(name string, def string, usage string) *string {
	sv := newStringValue(&def)
	set.value(sv, def, name, usage)
	return sv.ptr
}

// StringVar defines a string environment variable with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the environment variable.
func (set *VariableSet) StringVar(p *string, name string, def string, usage string) {
	if p == nil {
		panic("`p` cannot be nil")
	}
	sv := newStringValue(p)
	set.value(sv, def, name, usage)
}

// Int defines a int value with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the environment variable.
func (set *VariableSet) Int(name string, def int, usage string) *int {
	iv := newIntValue(&def)
	set.value(iv, def, name, usage)
	return iv.ptr
}

// IntVar defines a int environment variable with specified name, default value, and usage string.
// The argument p points to a int variable in which to store the value of the environment variable.
func (set *VariableSet) IntVar(p *int, name string, def int, usage string) {
	if p == nil {
		panic("`p` cannot be nil")
	}
	iv := newIntValue(p)
	set.value(iv, def, name, usage)
}

// Bool defines a bool value with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the environment variable.
func (set *VariableSet) Bool(name string, def bool, usage string) *bool {
	bv := newBoolValue(&def)
	set.value(bv, def, name, usage)
	return bv.ptr
}

// BoolVar defines a bool environment variable with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the environment variable.
func (set *VariableSet) BoolVar(p *bool, name string, def bool, usage string) {
	if p == nil {
		panic("`p` cannot be nil")
	}
	bv := newBoolValue(p)
	set.value(bv, def, name, usage)
}

// Parse parses the environment variables from argument. Must be called
// after all variables are defined and before variable are accessed by the program.
func (set *VariableSet) Parse(environments []string) error {

	keys := set.sortedKeys()
	for _, key := range keys {
		val := set.variables[key]

		osVal, found := set.lookup(environments, key)
		if !found {
			continue
		}

		err := val.value.Set(osVal)
		if err == nil {
			continue
		}

		e := fmt.Errorf("parse error %s: %w", key, err)
		switch set.errHandling {

		case ContinueOnError:
			return e

		case ExitOnError:
			w := set.output
			if w == nil {
				w = os.Stderr
			}
			set.fprintf(w, "%v\n", e)
			os.Exit(2)

		default:
			fallthrough
		case PanicOnError:
			panic(e)
		}
	}

	return nil
}

func (set *VariableSet) value(v Value, def interface{}, name string, usage string) {

	if set.variables == nil {
		set.variables = make(map[string]*variable)
	}

	_, ok := set.variables[name]
	if ok {
		msg := fmt.Sprintf("%s is already exist", name)
		panic(msg)
	}

	val := &variable{
		name:  name,
		value: v,
		def:   def,
		usage: usage,
	}
	set.variables[name] = val
}

func (set *VariableSet) defaultUsage() string {
	var sb strings.Builder

	if len(set.name) > 0 {
		set.fprintf(&sb, "Usage of %s\n", set.name)
	}

	keys := set.sortedKeys()
	for _, key := range keys {
		variable := set.variables[key]

		if variable.def != nil {
			set.fprintf(&sb, "  %s: default(%v)\n", variable.name, variable.def)
		} else {
			set.fprintf(&sb, "  %s:\n", variable.name)
		}

		for _, line := range strings.Split(variable.usage, "\n") {
			set.fprintf(&sb, "    %s\n", line)
		}
	}

	return sb.String()
}

func (_ *VariableSet) fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func (set *VariableSet) lookup(environments []string, key string) (string, bool) {

	for _, pair := range environments {

		k, v, err := set.splitKeyValue(pair)
		if err != nil {
			continue
		}

		if k == key {
			return v, true
		}
	}

	return "", false
}

func (_ *VariableSet) splitKeyValue(s string) (string, string, error) {
	index := strings.IndexByte(s, '=')
	if index < 0 {
		return "", "", fmt.Errorf("invalid format `%s`", s)
	}

	return s[:index], s[index+1:], nil
}

func (set *VariableSet) sortedKeys() []string {
	var keys []string
	for key := range set.variables {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

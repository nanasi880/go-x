package env

// ErrorHandling defines how VariableSet.Parse behaves if the parse fails.
type ErrorHandling int

// These constants cause VariableSet.Parse to behave as described if the parse fails.
const (
	ContinueOnError ErrorHandling = iota // Return a descriptive error.
	ExitOnError                          // Call os.Exit(2).
	PanicOnError                         // Call panic with a descriptive error.
)


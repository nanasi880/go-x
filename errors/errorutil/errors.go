package errorutil

import (
	"fmt"
	"runtime"
)

type object struct {
	message    string
	stacktrace []runtime.Frame
}

func (o *object) Error() string {
	return o.message
}

func New(s string, args ...interface{}) error {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	return &object{
		message:    s,
		stacktrace: stacktrace(),
	}
}

func stacktrace() []runtime.Frame {
	var pc [50]uintptr
	n := runtime.Callers(2, pc[:])
	if n <= 0 {
		return nil
	}

	var (
		frames     = runtime.CallersFrames(pc[:n])
		stacktrace = make([]runtime.Frame, 0)
		zeroFrame  runtime.Frame
	)
	for {
		frame, _ := frames.Next()
		if frame == zeroFrame {
			break
		}
		stacktrace = append(stacktrace, frame)
	}
	return stacktrace
}

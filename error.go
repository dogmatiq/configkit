package configkit

import (
	"fmt"
	"strings"
)

// Error is an error representing a fault in an application's
// configuration.
type Error string

func (e Error) Error() string {
	return string(e)
}

// Recover recovers from a configuration related panic.
//
// It is intended to be used in a defer statement. If the panic value is a
// Error, it is assigned to *err.
func Recover(err *error) {
	if err == nil {
		panic("err must be a non-nil pointer")
	}

	switch v := recover().(type) {
	case Error:
		*err = v
	case nil:
		return
	default:
		panic(v)
	}
}

// Panicf panics with a new Error.
func Panicf(f string, v ...interface{}) {
	panic(Errorf(f, v...))
}

// Errorf returns a new ValidationEerror.
func Errorf(f string, v ...interface{}) Error {
	m := fmt.Sprintf(f, v...)
	m = strings.Replace(m, "an command", "a command", -1)
	m = strings.Replace(m, "a event", "an event", -1)
	return Error(m)
}

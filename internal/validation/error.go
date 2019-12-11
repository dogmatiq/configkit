package validation

import (
	"fmt"
	"strings"
)

// Error is an error representing a fault in an entity's configuration.
type Error string

func (e Error) Error() string {
	return string(e)
}

// Panicf panics with a new Error.
func Panicf(f string, v ...interface{}) {
	panic(Errorf(f, v...))
}

// Errorf returns a new Error.
func Errorf(f string, v ...interface{}) Error {
	m := fmt.Sprintf(f, v...)
	m = strings.ReplaceAll(m, "an command", "a command")
	m = strings.ReplaceAll(m, "a event", "an event")
	m = strings.ReplaceAll(m, "an timeout", "a timeout")
	return Error(m)
}

package config

import (
	"fmt"
	"strings"
)

// catch calls fn(), and recovers from any panic where the value is a
// ConfigurationError by returning that error.
func catch(fn func()) (err error) {
	defer func() {
		switch r := recover().(type) {
		case Error:
			err = r
		case nil:
			return
		default:
			panic(r)
		}
	}()

	fn()

	return
}

// panicf panics with a new configuration error.
func panicf(f string, v ...interface{}) {
	panic(errorf(f, v...))
}

// errorf returns a new configuration error.
func errorf(f string, v ...interface{}) Error {
	m := fmt.Sprintf(f, v...)
	m = strings.Replace(m, "an command", "a command", -1)
	m = strings.Replace(m, "a event", "an event", -1)
	return Error(m)
}

// Error is an error representing a fault in an application's
// configuration.
type Error string

func (e Error) Error() string {
	return string(e)
}

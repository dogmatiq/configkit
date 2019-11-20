package configkit

import (
	"fmt"
	"strings"
)

// ValidationError is an error representing a fault in an application's
// configuration.
type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

// catch calls fn(), and recovers from any panic where the value is a
// ConfigurationError by returning that error.
func catch(fn func()) (err error) {
	defer func() {
		switch r := recover().(type) {
		case ValidationError:
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
func errorf(f string, v ...interface{}) ValidationError {
	m := fmt.Sprintf(f, v...)
	m = strings.Replace(m, "an command", "a command", -1)
	m = strings.Replace(m, "a event", "an event", -1)
	return ValidationError(m)
}

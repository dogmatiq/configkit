package configkit

import (
	"github.com/dogmatiq/configkit/internal/validation"
)

// Error is an error representing a fault in an entity's configuration.
type Error = validation.Error

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

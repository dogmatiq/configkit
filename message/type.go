package message

import (
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Type represents the type of a Dogma message.
type Type struct {
	n  Name
	rt reflect.Type
}

// TypeOf returns the message type of m.
func TypeOf(m dogma.Message) Type {
	rt := reflect.TypeOf(m)
	n := buildName(rt, true)

	return Type{
		Name{n},
		rt,
	}
}

// Name returns the fully-qualified name for the Go type.
//
// It panics if t.IsZero() returns true.
func (t Type) Name() Name {
	if t.IsZero() {
		panic("can not obtain name of zero-value type")
	}

	return t.n
}

// ReflectType returns the reflect.Type of the message.
//
// It panics if t.IsZero() returns true.
func (t Type) ReflectType() reflect.Type {
	if t.IsZero() {
		panic("can not obtain reflect type of zero-value type")
	}

	return t.rt
}

// String returns a human-readable name for the type.
//
// The returned name is not necessarily globally-unique.
func (t Type) String() string {
	return buildName(t.rt, false)
}

// IsZero returns true if t is the zero-value.
func (t Type) IsZero() bool {
	return t.n.IsZero()
}

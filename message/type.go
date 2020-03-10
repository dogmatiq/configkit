package message

import (
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename"
	"github.com/dogmatiq/dogma"
)

// TypeCollection is an interface for containers of message types.
type TypeCollection interface {
	// Has returns true if t is in the container.
	Has(t Type) bool

	// HasM returns true if TypeOf(m) is in the container.
	HasM(m dogma.Message) bool

	// Range invokes fn once for each type in the container.
	//
	// Iteration stops when fn returns false or once fn has been invoked for all
	// types in the container.
	//
	// It returns true if fn returned true for all types.
	Range(fn func(Type) bool) bool
}

// IsIntersectingT returns true if a and b are intersecting.
//
// That is, it returns true if a and b contain any of the same types.
//
// See https://en.wikipedia.org/wiki/Set_(mathematics)#Intersections.
func IsIntersectingT(a, b TypeCollection) bool {
	return !a.Range(func(t Type) bool {
		return !b.Has(t)
	})
}

// IsSubsetT returns true if sub is a (non-strict) subset of sup.
//
// That is, it returns true if sup contains all of the types in sub.
//
// See https://en.wikipedia.org/wiki/Set_(mathematics)#Subsets.
func IsSubsetT(sub, sup TypeCollection) bool {
	return sub.Range(func(t Type) bool {
		return sup.Has(t)
	})
}

// Type represents the type of a Dogma message.
type Type struct {
	n  Name
	rt reflect.Type
}

// TypeOf returns the message type of m.
func TypeOf(m dogma.Message) Type {
	if m == nil {
		panic("message must not be nil")
	}

	return TypeFromReflect(reflect.TypeOf(m))
}

// TypeFromReflect returns the message type of the given reflect type.
func TypeFromReflect(rt reflect.Type) Type {
	// This is a compile-time assertion that the dogma.Message interface is
	// empty. If this fails to compile, this function needs additional logic to
	// verify that the type represented by rt actually implements dogma.Message.
	var _ interface{} = (dogma.Message)(nil)

	n := typename.Of(rt)

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
	return t.rt.String()
}

// IsZero returns true if t is the zero-value.
func (t Type) IsZero() bool {
	return t.n.IsZero()
}

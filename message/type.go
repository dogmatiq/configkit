package message

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/dogma"
)

// TypeCollection is an interface for containers of message types.
type TypeCollection interface {
	// Has returns true if t is in the container.
	Has(t Type) bool

	// HasM returns true if TypeOf(m) is in the container.
	HasM(m dogma.Message) bool

	// Len returns the number of names in the collection.
	Len() int

	// Range invokes fn once for each type in the container.
	//
	// Iteration stops when fn returns false or once fn has been invoked for all
	// types in the container.
	//
	// It returns true if fn returned true for all types.
	Range(fn func(Type) bool) bool
}

// IsEqualSetT returns true if a and b are equal.
//
// That is, it returns true if and only if every element of a is an element of
// b, and vice-versa.
func IsEqualSetT(a, b TypeCollection) bool {
	return IsSubsetT(a, b) && IsSubsetT(b, a)
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
	r  Role
	rt reflect.Type
}

// TypeOf returns the message type of m.
func TypeOf(m dogma.Message) Type {
	if m == nil {
		panic("message must not be nil")
	}

	return TypeFromReflect(reflect.TypeOf(m))
}

// TypeFor returns the message type for T.
func TypeFor[T dogma.Message]() Type {
	return TypeFromReflect(reflect.TypeFor[T]())
}

var (
	interfaces = map[Role]reflect.Type{
		CommandRole: reflect.TypeFor[dogma.Command](),
		EventRole:   reflect.TypeFor[dogma.Event](),
		TimeoutRole: reflect.TypeFor[dogma.Timeout](),
	}
)

// TypeFromReflect returns the message type of the given reflect type.
func TypeFromReflect(rt reflect.Type) Type {
	n := goreflect.NameOf(rt)

	for r, i := range interfaces {
		if rt.Implements(i) {
			return Type{
				Name{n},
				r,
				rt,
			}
		}
	}

	panic(fmt.Sprintf(
		"%s does not implement dogma.Command, dogma.Event, or dogma.Timeout",
		rt,
	))
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

// Role returns the role of the message.
func (t Type) Role() Role {
	if t.IsZero() {
		panic("can not obtain role of zero-value type")
	}
	return t.r
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
	return typeToString(t.rt)
}

// IsZero returns true if t is the zero-value.
func (t Type) IsZero() bool {
	return t.n.IsZero()
}

func typeToString(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return "*" + typeToString(t.Elem())
	}

	str := t.String()

	if pkg := t.PkgPath(); pkg != "" {
		str = strings.ReplaceAll(str, pkg+".", "")
	}

	return str
}

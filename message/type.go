package message

import (
	"reflect"
	"strings"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/dogma"
)

// Type represents the type of a Dogma message.
type Type struct {
	n  Name
	k  Kind
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

// TypeFromReflect returns the message type of the given reflect type.
func TypeFromReflect(rt reflect.Type) Type {
	return Type{
		Name{goreflect.NameOf(rt)},
		kindFromReflect(rt),
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

// Kind returns the kind of the message represented by t.
//
// It panics of t does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func (t Type) Kind() Kind {
	if t.IsZero() {
		panic("can not obtain kind of zero-value type")
	}

	return t.k
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

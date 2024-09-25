package message

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"github.com/dogmatiq/dogma"
)

// Name is the fully-qualified name of a message type.
type Name struct {
	n string
}

// NameOf returns the fully-qualified type name of v.
func NameOf(m dogma.Message) Name {
	if m == nil {
		panic("message must not be nil")
	}

	rt := reflect.TypeOf(m)
	n := goreflect.NameOf(rt)
	return Name{n}
}

// NameFor returns the message name for T.
func NameFor[T dogma.Message]() Name {
	rt := reflect.TypeFor[T]()
	n := goreflect.NameOf(rt)
	return Name{n}
}

// NameFromStaticType returns the fully-qualified type name of t.
func NameFromStaticType(t types.Type) Name {
	if t == nil {
		panic("type must not be nil")
	}

	n := gotypes.NameOf(t)
	return Name{n}
}

// String returns the fully-qualified type name as a string.
func (n Name) String() string {
	return n.n
}

// IsZero returns true if n is the zero-value.
func (n Name) IsZero() bool {
	return n.n == ""
}

// MarshalText returns a UTF-8 representation of the name.
func (n Name) MarshalText() ([]byte, error) {
	if n.n == "" {
		return nil, fmt.Errorf("can not marshal empty name")
	}

	return []byte(n.n), nil
}

// UnmarshalText unmarshals a name from its UTF-8 representation.
func (n *Name) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return fmt.Errorf("can not unmarshal empty name")
	}

	n.n = string(text)
	return nil
}

// MarshalBinary returns a binary representation of the name.
func (n Name) MarshalBinary() ([]byte, error) {
	return n.MarshalText()
}

// UnmarshalBinary unmarshals a type from its binary representation.
func (n *Name) UnmarshalBinary(data []byte) error {
	return n.UnmarshalText(data)
}

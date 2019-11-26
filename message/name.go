package message

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename"
	"github.com/dogmatiq/dogma"
)

// NameCollection is an interface for containers of message names.
type NameCollection interface {
	// Has returns true if n is in the container.
	Has(n Name) bool

	// HasM returns true if NameOf(m) is in the container.
	HasM(m dogma.Message) bool

	// Each invokes fn once for each name in the container.
	//
	// Iteration stops when fn returns false or once fn has been invoked for all
	// names in the container.
	//
	// It returns true if fn returned true for all names.
	Each(fn func(Name) bool) bool
}

// Name is the fully-qualified name of a message type.
type Name struct {
	n string
}

// NameOf returns the fully-qualified type name of v.
func NameOf(v dogma.Message) Name {
	rt := reflect.TypeOf(v)
	n := typename.Of(rt)
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

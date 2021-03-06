package message

import (
	"fmt"
	"go/types"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"github.com/dogmatiq/dogma"
)

// NameCollection is an interface for containers of message names.
type NameCollection interface {
	// Has returns true if n is in the container.
	Has(n Name) bool

	// HasM returns true if NameOf(m) is in the container.
	HasM(m dogma.Message) bool

	// Len returns the number of names in the collection.
	Len() int

	// Range invokes fn once for each name in the container.
	//
	// Iteration stops when fn returns false or once fn has been invoked for all
	// names in the container.
	//
	// It returns true if fn returned true for all names.
	Range(fn func(Name) bool) bool
}

// IsEqualSetN returns true if a and b are equal.
//
// That is, it returns true if and only if every element of a is an element of
// b, and vice-versa.
func IsEqualSetN(a, b NameCollection) bool {
	return IsSubsetN(a, b) && IsSubsetN(b, a)
}

// IsIntersectingN returns true if a and b are intersecting.
//
// That is, it returns true if a and b contain any of the same names.
//
// See https://en.wikipedia.org/wiki/Set_(mathematics)#Intersections.
func IsIntersectingN(a, b NameCollection) bool {
	return !a.Range(func(n Name) bool {
		return !b.Has(n)
	})
}

// IsSubsetN returns true if sub is a (non-strict) subset of sup.
//
// That is, it returns true if sup contains all of the names in sub.
//
// See https://en.wikipedia.org/wiki/Set_(mathematics)#Subsets.
func IsSubsetN(sub, sup NameCollection) bool {
	return sub.Range(func(n Name) bool {
		return sup.Has(n)
	})
}

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

// NameFromType returns the fully-qualified type name of t.
func NameFromType(t types.Type) Name {
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

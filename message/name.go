package message

import (
	"fmt"
	"reflect"

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
	n := buildName(rt, true)
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

// buildName returns the name of the given type.
//
// If qual is true, it produces package-path-qualified names.
// It panics if rt does not implement dogma.Message.
func buildName(rt reflect.Type, qual bool) string {
	// This is a static assertion that the dogma.Message interface is empty. If
	// this fails to compile, additional logic must be added to verify that rt
	// truly does implement dogma.Message.
	var _ dogma.Message = (interface{})(nil)

	if rt.Name() != "" {
		return buildDefinedName(rt, qual)
	}

	switch rt.Kind() {
	case reflect.Ptr:
		return fmt.Sprintf("*%s", buildName(rt.Elem(), qual))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", buildName(rt.Elem(), qual))
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", rt.Len(), buildName(rt.Elem(), qual))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", buildName(rt.Key(), qual), buildName(rt.Elem(), qual))
	case reflect.Chan:
		return buildChanName(rt, qual)
	case reflect.Struct:
		return buildStructName(rt, qual)
	case reflect.Func:
		return buildFuncName(rt, qual)
	}

	// CODE COVERAGE: Type is likely an interface. This path is unreachable
	// using the public API.
	panic("unsupported type: " + rt.String())
}

func buildDefinedName(rt reflect.Type, qual bool) string {
	if !qual {
		return rt.String()
	}

	var name string

	if p := rt.PkgPath(); p != "" {
		name += p
		name += `.`
	}

	return name + rt.Name()
}

func buildChanName(rt reflect.Type, qual bool) string {
	elem := buildName(rt.Elem(), qual)

	switch rt.ChanDir() {
	case reflect.RecvDir: // <-chan
		return "<-chan " + elem
	case reflect.SendDir: // chan<-
		return "chan<- " + elem
	default:
		return "chan " + elem
	}
}

func buildStructName(rt reflect.Type, qual bool) string {
	name := "struct {"

	for i := 0; i < rt.NumField(); i++ {
		if i > 0 {
			name += "; "
		}

		f := rt.Field(i)
		name += f.Name

		if !f.Anonymous {
			name += " "
			name += buildName(f.Type, qual)
		}
	}

	name += "}"

	return name
}

func buildFuncName(rt reflect.Type, qual bool) string {
	name := "func("
	for i := 0; i < rt.NumIn(); i++ {
		if i > 0 {
			name += ", "
		}

		name += buildName(rt.In(i), qual)
	}
	name += ")"

	if n := rt.NumOut(); n != 0 {
		name += " "

		if n > 1 {
			name += "("
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				name += ", "
			}

			name += buildName(rt.Out(i), qual)
		}

		if n > 1 {
			name += ")"
		}
	}

	return name
}

package configkit

import "reflect"

// TypeName is a fully-qualified name for a Go type.
type TypeName string

// NewTypeName returns the fully-qualified name of the given type.
func NewTypeName(rt reflect.Type) TypeName {
	panic("not implemented")
}

// TypeNameOf returns the fully-qualified type name of the given value.
func TypeNameOf(v interface{}) TypeName {
	return NewTypeName(
		reflect.TypeOf(v),
	)
}

package gotypes

import (
	"fmt"
	"go/types"
)

// NameOf returns the fully-qualified name of the given type.
func NameOf(t types.Type) string {
	switch t := t.(type) {
	case *types.Basic:
		return t.String()
	case *types.Pointer:
		return t.String()
	case *types.Slice:
		return t.String()
	case *types.Array:
		return t.String()
	case *types.Map:
		return t.String()
	case *types.Chan:
		return buildChanName(t)
	case *types.Interface:
		return buildInterfaceName(t)
	case *types.Named:
		return t.String()
	case *types.Struct:
		return buildStructName(t)
	case *types.Signature:
		return "func" + buildFuncSignature(t)
	}

	// COVERAGE: This panic is only possible if an additional implementation of
	// the types.Type interface is introduced to represent a new first-class
	// type in Go. Not all implementations of types.Type represent first-class
	// types. For instance, types.Tuple can only be a part of signatures or
	// multiple assignments.
	panic(fmt.Sprintf("unknown type %s", t))
}

func buildChanName(c *types.Chan) string {
	elem := NameOf(c.Elem())

	switch c.Dir() {
	case types.RecvOnly: // <-chan
		return "<-chan " + elem
	case types.SendOnly: // chan<-
		return "chan<- " + elem
	default:
		return "chan " + elem
	}
}

func buildInterfaceName(t *types.Interface) string {
	name := "interface {"

	n := t.NumMethods()

	if n > 0 {
		name += " "

		for i := 0; i < n; i++ {
			if i > 0 {
				name += "; "
			}

			m := t.Method(i)
			name += m.Name()
			name += buildFuncSignature(m.Type().(*types.Signature))
		}

		name += " "
	}

	name += "}"

	return name
}

func buildStructName(s *types.Struct) string {
	name := "struct {"

	n := s.NumFields()

	if n > 0 {
		name += " "

		for i := 0; i < n; i++ {
			if i > 0 {
				name += "; "
			}

			f := s.Field(i)

			if !f.Anonymous() {
				name += f.Name()
				name += " "
			}

			name += NameOf(f.Type())
		}

		name += " "
	}

	name += "}"

	return name
}

func buildFuncSignature(s *types.Signature) string {
	name := "("
	for i := 0; i < s.Params().Len(); i++ {
		if i > 0 {
			name += ", "
		}

		name += NameOf(s.Params().At(i).Type())
	}
	name += ")"

	if n := s.Results().Len(); n != 0 {
		name += " "

		if n > 1 {
			name += "("
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				name += ", "
			}

			name += NameOf(s.Results().At(i).Type())
		}

		if n > 1 {
			name += ")"
		}
	}

	return name
}

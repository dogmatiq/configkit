package gotypes

import (
	"go/types"
)

// Of returns the fully-qualified name of the given type.
func Of(t types.Type) string {
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
	default:
		return t.String()
	}
}

func buildChanName(c *types.Chan) string {
	elem := Of(c.Elem())

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

			name += Of(f.Type())
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

		name += Of(s.Params().At(i).Type())
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

			name += Of(s.Results().At(i).Type())
		}

		if n > 1 {
			name += ")"
		}
	}

	return name
}

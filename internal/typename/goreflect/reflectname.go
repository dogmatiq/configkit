package goreflect

import (
	"fmt"
	"reflect"
)

// NameOf returns the fully-qualified name of the given type.
//
// The return value is similar to the string returned by ReflectType.String()
// but shows fully-qualified package paths, not just package names.
func NameOf(rt reflect.Type) string {
	if rt.Name() != "" {
		return buildDefinedName(rt)
	}

	// only the "composite" types can be unnabled:
	switch rt.Kind() {
	case reflect.Ptr:
		return fmt.Sprintf("*%s", NameOf(rt.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", NameOf(rt.Elem()))
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", rt.Len(), NameOf(rt.Elem()))
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", NameOf(rt.Key()), NameOf(rt.Elem()))
	case reflect.Chan:
		return buildChanName(rt)
	case reflect.Interface:
		return buildInterfaceName(rt)
	case reflect.Struct:
		return buildStructName(rt)
	default: // reflect.Func
		return "func" + buildFuncSignature(rt)
	}
}

func buildDefinedName(rt reflect.Type) string {
	var name string

	if p := rt.PkgPath(); p != "" {
		name += p
		name += `.`
	}

	return name + rt.Name()
}

func buildChanName(rt reflect.Type) string {
	elem := NameOf(rt.Elem())

	switch rt.ChanDir() {
	case reflect.RecvDir: // <-chan
		return "<-chan " + elem
	case reflect.SendDir: // chan<-
		return "chan<- " + elem
	default:
		return "chan " + elem
	}
}

func buildInterfaceName(rt reflect.Type) string {
	name := "interface {"

	n := rt.NumMethod()

	if n > 0 {
		name += " "

		for i := 0; i < n; i++ {
			if i > 0 {
				name += "; "
			}

			m := rt.Method(i)

			name += m.Name
			name += buildFuncSignature(m.Type)
		}

		name += " "
	}

	name += "}"

	return name
}

func buildStructName(rt reflect.Type) string {
	name := "struct {"

	n := rt.NumField()

	if n > 0 {
		name += " "

		for i := 0; i < n; i++ {
			if i > 0 {
				name += "; "
			}

			f := rt.Field(i)

			if !f.Anonymous {
				name += f.Name
				name += " "
			}

			name += NameOf(f.Type)
		}

		name += " "
	}

	name += "}"

	return name
}

func buildFuncSignature(rt reflect.Type) string {
	name := "("
	for i := 0; i < rt.NumIn(); i++ {
		if i > 0 {
			name += ", "
		}

		name += NameOf(rt.In(i))
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

			name += NameOf(rt.Out(i))
		}

		if n > 1 {
			name += ")"
		}
	}

	return name
}

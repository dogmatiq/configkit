package goreflect_test

import (
	"reflect"

	. "github.com/dogmatiq/configkit/internal/typename"
	. "github.com/dogmatiq/configkit/internal/typename/goreflect"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("func Of()", func() {
	Declare(
		"../testdata",
		func(file string) string {
			switch file {
			case "basic":
				return Of(reflect.TypeOf(Basic))
			case "named":
				return Of(reflect.TypeOf(Named))
			case "pointer":
				return Of(reflect.TypeOf(Pointer))
			case "slice":
				return Of(reflect.TypeOf(Slice))
			case "array":
				return Of(reflect.TypeOf(Array))
			case "map":
				return Of(reflect.TypeOf(Map))
			case "channel":
				return Of(reflect.TypeOf(Channel))
			case "receive-only-channel":
				return Of(reflect.TypeOf(ReceiveOnlyChannel))
			case "send-only-channel":
				return Of(reflect.TypeOf(SendOnlyChannel))
			case "iface":
				return Of(reflect.TypeOf(&Interface).Elem())
			case "iface-with-single-method":
				return Of(reflect.TypeOf(&InterfaceWithSingleMethod).Elem())
			case "iface-with-multiple-methods":
				return Of(reflect.TypeOf(&InterfaceWithMultipleMethods).Elem())
			case "struct":
				return Of(reflect.TypeOf(Struct))
			case "struct-with-single-field":
				return Of(reflect.TypeOf(StructWithSingleField))
			case "struct-with-multiple-fields":
				return Of(reflect.TypeOf(StructWithMultipleFields))
			case "struct-with-embedded-type":
				return Of(reflect.TypeOf(StructWithEmbeddedType))
			case "nullary":
				return Of(reflect.TypeOf(Nullary))
			case "func-with-multiple-input-params":
				return Of(reflect.TypeOf(FuncWithMultipleInputParams))
			case "func-with-single-output-param":
				return Of(reflect.TypeOf(FuncWithSingleOutputParam))
			case "func-with-multiple-output-params":
				return Of(reflect.TypeOf(FuncWithMultipleOutputParams))
			default:
				Fail("no type is available for testing with '%s' file")
			}

			return ""
		})
})

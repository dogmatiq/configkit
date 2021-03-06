package gotypes_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"

	. "github.com/dogmatiq/configkit/internal/typename/gotypes"
	. "github.com/dogmatiq/configkit/internal/typename/internal/typenametest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func NameOf()", func() {
	var varmap map[string]types.Type

	BeforeSuite(func() {
		bb, err := ioutil.ReadFile("../internal/typenametest/types.go")
		Expect(err).ShouldNot(HaveOccurred())

		fset := token.NewFileSet()

		f, err := parser.ParseFile(
			fset,
			"types.go",
			string(bb),
			parser.ParseComments,
		)
		Expect(err).ShouldNot(HaveOccurred())

		conf := types.Config{Importer: importer.Default()}

		pkg, err := conf.Check(
			"github.com/dogmatiq/configkit/internal/typename/internal/typenametest",
			fset,
			[]*ast.File{f},
			nil,
		)
		Expect(err).ShouldNot(HaveOccurred())

		names := pkg.Scope().Names()

		varmap = make(map[string]types.Type, len(names))

		for _, name := range names {
			if obj := pkg.Scope().Lookup(name); obj != nil {
				varmap[obj.Name()] = obj.Type()
			}
		}
	})

	Declare(
		func(file string) string {
			switch file {
			case "basic":
				return NameOf(varmap["Basic"])
			case "named":
				return NameOf(varmap["Named"])
			case "pointer":
				return NameOf(varmap["Pointer"])
			case "slice":
				return NameOf(varmap["Slice"])
			case "array":
				return NameOf(varmap["Array"])
			case "map":
				return NameOf(varmap["Map"])
			case "channel":
				return NameOf(varmap["Channel"])
			case "receive-only-channel":
				return NameOf(varmap["ReceiveOnlyChannel"])
			case "send-only-channel":
				return NameOf(varmap["SendOnlyChannel"])
			case "iface":
				return NameOf(varmap["Interface"])
			case "iface-with-single-method":
				return NameOf(varmap["InterfaceWithSingleMethod"])
			case "iface-with-multiple-methods":
				return NameOf(varmap["InterfaceWithMultipleMethods"])
			case "struct":
				return NameOf(varmap["Struct"])
			case "struct-with-single-field":
				return NameOf(varmap["StructWithSingleField"])
			case "struct-with-multiple-fields":
				return NameOf(varmap["StructWithMultipleFields"])
			case "struct-with-embedded-type":
				return NameOf(varmap["StructWithEmbeddedType"])
			case "nullary":
				return NameOf(varmap["Nullary"])
			case "func-with-multiple-input-params":
				return NameOf(varmap["FuncWithMultipleInputParams"])
			case "func-with-single-output-param":
				return NameOf(varmap["FuncWithSingleOutputParam"])
			case "func-with-multiple-output-params":
				return NameOf(varmap["FuncWithMultipleOutputParams"])
			}

			Fail("no type is available for testing with '%s' file")
			return ""
		})
})

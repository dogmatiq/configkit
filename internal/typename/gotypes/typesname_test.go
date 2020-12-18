package gotypes_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"

	. "github.com/dogmatiq/configkit/internal/typename"
	. "github.com/dogmatiq/configkit/internal/typename/gotypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var _ = Describe("func Of()", func() {
	var varmap map[string]types.Type

	BeforeSuite(func() {
		bb, err := ioutil.ReadFile("../types.go")
		Expect(err).ShouldNot(HaveOccurred())

		fset := token.NewFileSet()

		f, err := parser.ParseFile(
			fset,
			"types.go",
			string(bb),
			parser.ParseComments,
		)
		Expect(err).ShouldNot(HaveOccurred())
		files := []*ast.File{f}

		pkg := types.NewPackage(
			"github.com/dogmatiq/configkit/internal/typename",
			"typename",
		)

		built, _, err := ssautil.BuildPackage(
			&types.Config{
				Importer: importer.Default(),
			},
			fset,
			pkg,
			files,
			ssa.SanityCheckFunctions,
		)
		if err != nil {
			panic(err)
		}

		varmap = make(map[string]types.Type, len(built.Members))

		for _, m := range built.Members {
			if g, ok := m.(*ssa.Global); ok {
				if g.Object() != nil {
					varmap[g.Name()] = g.Object().Type()
				}
			}
		}
	})

	Declare(
		"../testdata",
		func(file string) string {
			switch file {
			case "basic":
				return Of(varmap["Basic"])
			case "named":
				return Of(varmap["Named"])
			case "pointer":
				return Of(varmap["Pointer"])
			case "slice":
				return Of(varmap["Slice"])
			case "array":
				return Of(varmap["Array"])
			case "map":
				return Of(varmap["Map"])
			case "channel":
				return Of(varmap["Channel"])
			case "receive-only-channel":
				return Of(varmap["ReceiveOnlyChannel"])
			case "send-only-channel":
				return Of(varmap["SendOnlyChannel"])
			case "iface":
				return Of(varmap["Interface"])
			case "iface-with-single-method":
				return Of(varmap["InterfaceWithSingleMethod"])
			case "iface-with-multiple-methods":
				return Of(varmap["InterfaceWithMultipleMethods"])
			case "struct":
				return Of(varmap["Struct"])
			case "struct-with-single-field":
				return Of(varmap["StructWithSingleField"])
			case "struct-with-multiple-fields":
				return Of(varmap["StructWithMultipleFields"])
			case "struct-with-embedded-type":
				return Of(varmap["StructWithEmbeddedType"])
			case "nullary":
				return Of(varmap["Nullary"])
			case "func-with-multiple-input-params":
				return Of(varmap["FuncWithMultipleInputParams"])
			case "func-with-single-output-param":
				return Of(varmap["FuncWithSingleOutputParam"])
			case "func-with-multiple-output-params":
				return Of(varmap["FuncWithMultipleOutputParams"])
			default:
				Fail("no type is available for testing with '%s' file")
			}

			return ""
		})
})

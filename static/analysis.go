package static

import (
	"go/constant"
	"go/types"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"golang.org/x/tools/go/ssa"
)

// analyzeApplication analyzes a type that implements the dogma.Application
// interface to deduce it's application configuration.
func analyzeApplication(pkg *ssa.Package, typ types.Type) configkit.Application {
	app := &entity.Application{
		TypeNameValue: gotypes.NameOf(typ),
	}

	fn := pkg.Prog.LookupMethod(
		typ,
		pkg.Pkg,
		"Configure",
	)

	for _, c := range findConfigurerCalls(fn) {
		switch c.Common().Method.Name() {
		case "Identity":
			app.IdentityValue = analyzeIdentityCall(c)
		}
	}

	return app
}

// findConfigurerCalls returns all of the calls to methods on the Dogma
// application or handler "configurer" within the given function.
func findConfigurerCalls(fn *ssa.Function) []*ssa.Call {
	// Since fn is expected to be a method, the first parameter is the receiver,
	// so the second parameter is the configurer itself.
	configurer := fn.Params[1]

	var calls []*ssa.Call

	for _, b := range fn.Blocks {
		for _, i := range b.Instrs {
			if c, ok := i.(*ssa.Call); ok &&
				c.Common().Value.Name() == configurer.Name() {
				calls = append(calls, c)
			}
		}
	}

	return calls
}

// analyzeIdentityCall analyzes the arguments in a call to a configurer's
// Identity() method to determine the application's or handler's identity.
func analyzeIdentityCall(c *ssa.Call) configkit.Identity {
	var ident configkit.Identity

	if c1, ok := c.Call.Args[0].(*ssa.Const); ok {
		ident.Name = constant.StringVal(c1.Value)
	}

	if c2, ok := c.Call.Args[1].(*ssa.Const); ok {
		ident.Key = constant.StringVal(c2.Value)
	}

	return ident
}

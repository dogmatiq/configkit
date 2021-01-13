package static

import (
	"go/constant"
	"go/types"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"golang.org/x/tools/go/ssa"
)

// analyzeApplication analyzes the member of the package and populates the
// collected datainto entity.Application.
//
// The second parameter is the type that actually implements
// github.com/dogmatiq/dogma.Application interface.
func analyzeApplication(m ssa.Member, typ types.Type) *entity.Application {
	app := &entity.Application{
		TypeNameValue: gotypes.NameOf(typ),
	}

	fn := m.Package().Prog.LookupMethod(
		typ,
		m.Package().Pkg,
		"Configure",
	)

	for _, c := range findConfigurerCalls(fn) {
		switch c.Common().Method.Name() {
		case "Identity":
			app.IdentityValue = parseIdentity(c)
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

// parseIdentity parses arguments passed to the .Identity() method.
func parseIdentity(c *ssa.Call) (ident configkit.Identity) {
	if c1, ok := c.Call.Args[0].(*ssa.Const); ok {
		ident.Name = constant.StringVal(c1.Value)
	}

	if c2, ok := c.Call.Args[1].(*ssa.Const); ok {
		ident.Key = constant.StringVal(c2.Value)
	}

	return
}

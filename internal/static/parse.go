package static

import (
	"go/constant"
	"go/types"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"golang.org/x/tools/go/ssa"
)

// parse parses the member of the package and populates it into
// entity.Application.
//
// The second parameter is the type that actually implements
// github.com/dogmatiq/dogma.Application interface.
func parse(m ssa.Member, typ types.Type) *entity.Application {
	app := &entity.Application{
		TypeNameValue: gotypes.NameOf(typ),
	}

	fn := m.Package().Prog.LookupMethod(
		typ,
		m.Package().Pkg,
		"Configure",
	)

	// Since `.Configure()` is the method, the first parameter is the receiver,
	// so we need to extract the second paramenter as
	// `dogma.ApplicationConfigurer`
	configurer := fn.Params[1]

	for _, c := range filterCalls(fn) {
		// Check if we are calling the methods on the configurer.
		if c.Common().Value.Name() == configurer.Name() {
			switch c.Common().Method.Name() {
			case "Identity":
				app.IdentityValue = parseIdentity(c)
			}
		}
	}

	return app
}

// filterCalls filters and returns all function/method call instructions within
// the given function.
func filterCalls(fn *ssa.Function) (calls []*ssa.Call) {
	for _, b := range fn.Blocks {
		for _, i := range b.Instrs {
			if c, ok := i.(*ssa.Call); ok {
				calls = append(calls, c)
			}
		}
	}

	return
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

package static

import (
	"fmt"
	"go/constant"
	"go/types"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"github.com/dogmatiq/configkit/message"
	"golang.org/x/tools/go/ssa"
)

// analyzeApplication analyzes a type that implements the dogma.Application
// interface to deduce its application configuration.
func analyzeApplication(prog *ssa.Program, typ types.Type) configkit.Application {
	app := &entity.Application{
		TypeNameValue:     gotypes.NameOf(typ),
		HandlersValue:     configkit.HandlerSet{},
		MessageNamesValue: configkit.EntityMessageNames{},
	}

	pkg := pkgOfNamedType(typ)
	fn := prog.LookupMethod(typ, pkg, "Configure")

	for _, c := range findConfigurerCalls(fn) {
		args := c.Common().Args

		switch c.Common().Method.Name() {
		case "Identity":
			app.IdentityValue = analyzeIdentityCall(c)
		case "RegisterAggregate":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				app.HandlersValue.Add(
					analyzeHandler(
						prog,
						mi.X.Type(),
						configkit.AggregateHandlerType,
					),
				)
			}
		case "RegisterProcess":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				app.HandlersValue.Add(
					analyzeHandler(
						prog,
						mi.X.Type(),
						configkit.ProcessHandlerType,
					),
				)
			}
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

// analyzeHandler analyzes a type that implements one of the Dogma handler
// interfaces to deduce its handler configuration.
func analyzeHandler(
	prog *ssa.Program,
	typ types.Type,
	ht configkit.HandlerType,
) configkit.Handler {
	hdr := &entity.Handler{
		HandlerTypeValue: ht,
		TypeNameValue:    typ.String(),
		MessageNamesValue: configkit.EntityMessageNames{
			Produced: message.NameRoles{},
			Consumed: message.NameRoles{},
		},
	}

	pkg := pkgOfNamedType(typ)
	fn := prog.LookupMethod(typ, pkg, "Configure")

	for _, c := range findConfigurerCalls(fn) {
		args := c.Common().Args

		switch c.Common().Method.Name() {
		case "Identity":
			hdr.IdentityValue = analyzeIdentityCall(c)
		case "ConsumesCommandType":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				hdr.MessageNamesValue.Consumed.Add(
					message.NameFromType(
						mi.X.Type(),
					),
					message.CommandRole,
				)
			}
		case "ConsumesEventType":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				hdr.MessageNamesValue.Consumed.Add(
					message.NameFromType(
						mi.X.Type(),
					),
					message.EventRole,
				)
			}
		case "ProducesCommandType":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				hdr.MessageNamesValue.Produced.Add(
					message.NameFromType(
						mi.X.Type(),
					),
					message.CommandRole,
				)
			}
		case "ProducesEventType":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				hdr.MessageNamesValue.Produced.Add(
					message.NameFromType(
						mi.X.Type(),
					),
					message.EventRole,
				)
			}
		case "SchedulesTimeoutType":
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				hdr.MessageNamesValue.Produced.Add(
					message.NameFromType(
						mi.X.Type(),
					),
					message.TimeoutRole,
				)
			}
			if mi, ok := args[0].(*ssa.MakeInterface); ok {
				hdr.MessageNamesValue.Consumed.Add(
					message.NameFromType(
						mi.X.Type(),
					),
					message.TimeoutRole,
				)
			}
		}
	}

	return hdr
}

// pkgOfNamedType returns the package in which typ is declared.
//
// It panics if typ is not a named type or a pointer to a named type.
func pkgOfNamedType(typ types.Type) *types.Package {
	switch t := typ.(type) {
	case *types.Named:
		return t.Obj().Pkg()
	case *types.Pointer:
		return pkgOfNamedType(t.Elem())
	default:
		panic(fmt.Sprintf("cannot determine package for anonymous or built-in type %v", typ))
	}
}

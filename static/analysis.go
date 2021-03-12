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

	for _, c := range findConfigurerCalls(prog, fn) {
		args := c.Common().Args

		switch c.Common().Method.Name() {
		case "Identity":
			app.IdentityValue = analyzeIdentityCall(c)
		case "RegisterAggregate":
			addHandlerFromArguments(
				prog,
				args,
				app.HandlersValue,
				configkit.AggregateHandlerType,
			)
		case "RegisterProcess":
			addHandlerFromArguments(
				prog,
				args,
				app.HandlersValue,
				configkit.ProcessHandlerType,
			)
		case "RegisterProjection":
			addHandlerFromArguments(
				prog,
				args,
				app.HandlersValue,
				configkit.ProjectionHandlerType,
			)
		}
	}

	return app
}

// findConfigurerCalls returns all of the calls to methods on the Dogma
// application or handler "configurer" within the given function.
//
// indices refers to the positions arguments that are the configurer. If none
// are provided it defaults to [1]. This accounts for the most common case where
// fn is the Configure() method on an application or handler. In this case the
// first parameter is the receiver, so the second parameter is the configurer
// itself.
func findConfigurerCalls(
	prog *ssa.Program,
	fn *ssa.Function,
	indices ...int,
) []*ssa.Call {
	if len(indices) == 0 {
		indices = []int{1}
	}

	configurers := map[ssa.Value]struct{}{}
	for _, i := range indices {
		v := fn.Params[i]
		configurers[v] = struct{}{}
	}

	var calls []*ssa.Call

	for _, b := range fn.Blocks {
		for _, i := range b.Instrs {
			if c, ok := i.(*ssa.Call); ok {
				if _, ok := configurers[c.Common().Value]; ok {
					// We've found a direct call to a method on the configurer.
					calls = append(calls, c)
				} else {
					// We've found a call to some other function or method. We
					// need to analyse the instructions within *that* function
					// to see if *it* makes any calls to the configurer.
					calls = append(
						calls,
						findIndirectConfigurerCalls(prog, configurers, c)...,
					)
				}
			}

		}
	}

	return calls
}

// findIndirectConfigurerCalls returns all of the calls to methods on the Dogma
// application or handler "configurer" within the an arbitrary function that has
// been called within a Configure() method.
func findIndirectConfigurerCalls(
	prog *ssa.Program,
	configurers map[ssa.Value]struct{},
	c *ssa.Call,
) []*ssa.Call {
	com := c.Common()

	var indices []int
	for i, arg := range com.Args {
		if _, ok := configurers[arg]; ok {
			indices = append(indices, i)
		}
	}

	if len(indices) == 0 {
		return nil
	}

	return findConfigurerCalls(
		prog,
		com.StaticCallee(),
		indices...,
	)
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

// addHandlerFromArguments analyzes args to deduce the configuration of a
// handler of type ht. It assumes that the handler is the first argument.
//
// If the first argument is not a pointer to ssa.MakeInterface instruction, this
// function has no effect; otherwise the handler is added to hs.
func addHandlerFromArguments(
	prog *ssa.Program,
	args []ssa.Value,
	hs configkit.HandlerSet,
	ht configkit.HandlerType,
) {
	if mi, ok := args[0].(*ssa.MakeInterface); ok {
		typ := mi.X.Type()

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

		for _, c := range findConfigurerCalls(prog, fn) {
			args := c.Common().Args

			switch c.Common().Method.Name() {
			case "Identity":
				hdr.IdentityValue = analyzeIdentityCall(c)
			case "ConsumesCommandType":
				addMessageFromArguments(
					args,
					hdr.MessageNamesValue.Consumed,
					message.CommandRole,
				)
			case "ConsumesEventType":
				addMessageFromArguments(
					args,
					hdr.MessageNamesValue.Consumed,
					message.EventRole,
				)
			case "ProducesCommandType":
				addMessageFromArguments(
					args,
					hdr.MessageNamesValue.Produced,
					message.CommandRole,
				)
			case "ProducesEventType":
				addMessageFromArguments(
					args,
					hdr.MessageNamesValue.Produced,
					message.EventRole,
				)
			case "SchedulesTimeoutType":
				addMessageFromArguments(
					args,
					hdr.MessageNamesValue.Consumed,
					message.TimeoutRole,
				)
				addMessageFromArguments(
					args,
					hdr.MessageNamesValue.Produced,
					message.TimeoutRole,
				)
			}
		}

		hs.Add(hdr)
	}
}

// addMessageFromArguments analyzes args to deduce the type of a message.
// It assumes that the message is always the first argument.
//
// If the first argument is not a pointer to ssa.MakeInterface instruction, this
// function has no effect; otherwise the message type is added to nr using the
// role given by r.
func addMessageFromArguments(
	args []ssa.Value,
	nr message.NameRoles,
	r message.Role,
) {
	if mi, ok := args[0].(*ssa.MakeInterface); ok {
		nr.Add(
			message.NameFromType(mi.X.Type()),
			r,
		)
	}
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

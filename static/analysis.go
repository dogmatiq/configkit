package static

import (
	"fmt"
	"go/constant"
	"go/types"
	"strings"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/typename/gotypes"
	"github.com/dogmatiq/configkit/message"
	"golang.org/x/tools/go/ssa"
)

// analyzeApplication analyzes a type that implements the dogma.Application
// interface to deduce its application configuration.
func analyzeApplication(
	prog *ssa.Program,
	dogmaPkg dogmaPackage,
	typ types.Type,
) configkit.Application {
	app := &application{
		TypeNameValue:     gotypes.NameOf(typ),
		HandlersValue:     configkit.HandlerSet{},
		MessageNamesValue: configkit.EntityMessages[message.Name]{},
	}

	pkg := pkgOfNamedType(typ)
	fn := prog.LookupMethod(typ, pkg, "Configure")

	for _, c := range findConfigurerCalls(prog, fn, make(map[string]types.Type)) {
		args := c.Common().Args

		switch c.Common().Method.Name() {
		case "Identity":
			app.IdentityValue = analyzeIdentityCall(c)
		case "RegisterAggregate":
			addHandlerFromArguments(
				prog,
				dogmaPkg,
				dogmaPkg.AggregateMessageHandler,
				args,
				app.HandlersValue,
				configkit.AggregateHandlerType,
			)
		case "RegisterProcess":
			addHandlerFromArguments(
				prog,
				dogmaPkg,
				dogmaPkg.ProcessMessageHandler,
				args,
				app.HandlersValue,
				configkit.ProcessHandlerType,
			)
		case "RegisterProjection":
			addHandlerFromArguments(
				prog,
				dogmaPkg,
				dogmaPkg.ProjectionMessageHandler,
				args,
				app.HandlersValue,
				configkit.ProjectionHandlerType,
			)
		case "RegisterIntegration":
			addHandlerFromArguments(
				prog,
				dogmaPkg,
				dogmaPkg.IntegrationMessageHandler,
				args,
				app.HandlersValue,
				configkit.IntegrationHandlerType,
			)
		}
	}

	app.MessageNamesValue = app.HandlersValue.MessageNames()

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
//
// The instantiatedTypes map is used to store the types that have been
// instantiated in the function. This is necessary because the SSA
// representation of a function does not include type information for the
// arguments, so we need to track this information ourselves. The keys are the
// names of the type parameters and the values are the concrete types that have
// been instantiated.
func findConfigurerCalls(
	prog *ssa.Program,
	fn *ssa.Function,
	instantiatedTypes map[string]types.Type,
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

	if fn.Synthetic != "" {
		populateInstantiatedTypes(fn, instantiatedTypes)
	}

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
						findIndirectConfigurerCalls(
							prog,
							configurers,
							c,
							instantiatedTypes,
						)...,
					)
				}
			}
		}
	}

	return calls
}

// populateInstantiatedTypes populates the instantiatedTypes map with the types
// that have been instantiated in the synthetic function.
func populateInstantiatedTypes(
	fn *ssa.Function,
	instantiatedTypes map[string]types.Type,
) {
	for _, b := range fn.Blocks {
		for _, i := range b.Instrs {
			if c, ok := i.(*ssa.ChangeType); ok {
				var (
					tpl *types.TypeParamList
					tal *types.TypeList
				)

				if ok, nt := namedType(c.Type()); ok {
					tpl = nt.TypeParams()
				}

				if ok, nt := namedType(c.X.Type()); ok {
					tal = nt.TypeArgs()
				}

				for i := 0; i < tpl.Len(); i++ {
					instantiatedTypes[tpl.At(i).String()] = tal.At(i)
				}
			}
		}
	}
}

// namedType returns true and the named type if t is a named type or a pointer
// to a named type.
func namedType(t types.Type) (ok bool, nt *types.Named) {
	switch t := t.(type) {
	case *types.Named:
		return true, t
	case *types.Pointer:
		return namedType(t.Elem())
	default:
		return false, nil
	}
}

// findIndirectConfigurerCalls returns all of the calls to methods on the Dogma
// application or handler "configurer" within the an arbitrary function that has
// been called within a Configure() method.
func findIndirectConfigurerCalls(
	prog *ssa.Program,
	configurers map[ssa.Value]struct{},
	c *ssa.Call,
	instantiatedTypes map[string]types.Type,
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

	if com.IsInvoke() {
		t, ok := instantiatedTypes[com.Value.Type().String()]
		if !ok {
			// If we cannot find any instantiated types in mapping, most likely
			// we hit the interface method and cannot analyze any further.
			return nil
		}

		return findConfigurerCalls(
			prog,
			prog.LookupMethod(t, com.Method.Pkg(), com.Method.Name()),
			instantiatedTypes,
			// don't pass indices here, as we are already in the method.
		)
	}

	return findConfigurerCalls(
		prog,
		com.StaticCallee(),
		instantiatedTypes,
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
func addHandlerFromArguments(
	prog *ssa.Program,
	dogmaPkg dogmaPackage,
	iface *types.Interface,
	args []ssa.Value,
	hs configkit.HandlerSet,
	ht configkit.HandlerType,
) {
	switch arg := args[0].(type) {
	case *ssa.MakeInterface:
		addHandlerFromType(prog, arg.X.Type(), hs, ht)
	case *ssa.Call:
		addHandlersFromAdaptorFunc(prog, dogmaPkg, iface, arg, hs, ht)
	}
}

// addHandlerFromType analyzes a type to deduce the configuration of a handler
// of type ht.
//
// The handler configuration is added to hs.
func addHandlerFromType(
	prog *ssa.Program,
	typ types.Type,
	hs configkit.HandlerSet,
	ht configkit.HandlerType,
) {
	pkg := pkgOfNamedType(typ)
	method := prog.LookupMethod(typ, pkg, "Configure")
	addHandlerFromConfigureMethod(prog, typ.String(), method, hs, ht)
}

// addHandlersFromAdaptorFunc analyzes the arguments of an "adaptor function" to
// deduce the configuration of handlers of type ht.
//
// The handler configurations are added to hs.
//
// Any function arguments that implement a Dogma style Configure() method are
// treated as a full handler implementation, even if they do not implement the
// other methods of the handler interface.
func addHandlersFromAdaptorFunc(
	prog *ssa.Program,
	dogmaPkg dogmaPackage,
	iface *types.Interface,
	call *ssa.Call,
	hs configkit.HandlerSet,
	ht configkit.HandlerType,
) {
	configureSig := prog.
		MethodSets.
		MethodSet(iface).
		Lookup(dogmaPkg.Package, "Configure").
		Type().(*types.Signature)

	for _, arg := range call.Call.Args {
		// If the argument to the adaptor function is an interface, use the
		// concrete type instead.
		if mi, ok := arg.(*ssa.MakeInterface); ok {
			arg = mi.X
		}

		typ := arg.Type()

		pkg, ok := tryPkgOfNamedType(typ)
		if !ok {
			// Argument is not a named type.
			continue
		}

		sel := prog.MethodSets.MethodSet(typ).Lookup(pkg, "Configure")
		if sel == nil {
			// Argument has no Configure() method.
			continue
		}

		method := prog.MethodValue(sel)
		if method == nil {
			// Configure() method exists, but it's "abstract".
			continue
		}

		if !types.Identical(method.Signature, configureSig) {
			// Configure() method exists, but it has a different signature than
			// is expected for this handler type.
			continue
		}

		addHandlerFromConfigureMethod(prog, typ.String(), method, hs, ht)
	}
}

// addHandlerFromConfigureMethod analyzes the body of a Configure() method to
// deduce the configuration of a handler of type ht.
//
// The handler configuration is added to hs.
func addHandlerFromConfigureMethod(
	prog *ssa.Program,
	handlerType string,
	method *ssa.Function,
	hs configkit.HandlerSet,
	ht configkit.HandlerType,
) {
	h := &handler{
		HandlerTypeValue:  ht,
		TypeNameValue:     handlerType,
		MessageNamesValue: configkit.EntityMessages[message.Name]{},
	}

	for _, c := range findConfigurerCalls(prog, method, make(map[string]types.Type)) {
		args := c.Common().Args

		switch c.Common().Method.Name() {
		case "Identity":
			h.IdentityValue = analyzeIdentityCall(c)
		case "Routes":
			addMessagesFromRoutes(h.MessageNamesValue, args)
		}
	}

	hs.Add(h)
}

// addMessagesFromRoutes analyzes the arguments in a call to a configurer's
// Routes() method to populate the messages that are produced and consumed by
// the handler.
func addMessagesFromRoutes(
	messages configkit.EntityMessages[message.Name],
	args []ssa.Value,
) {
	var mii []*ssa.MakeInterface
	for _, arg := range args {
		recurseSSAValues(
			arg,
			&[]ssa.Value{},
			func(v ssa.Value) bool {
				if v, ok := v.(*ssa.Call); ok {
					// We don't want to recurse into the call to Routes() method
					// itself.
					if v.Common().IsInvoke() &&
						v.Common().Method.Name() == "Routes" {
						return true
					}
				}

				// We do want to collect all of the MakeInterface instructions
				// that can potentially indicate boxing into Dogma route
				// interfaces.
				if v, ok := v.(*ssa.MakeInterface); ok {
					mii = append(mii, v)
					return true
				}

				return false
			},
		)
	}

	for _, mi := range mii {
		// If this is the boxing to the following interfaces,
		// we need to analyze the concrete types:
		switch mi.X.Type().String() {
		case "github.com/dogmatiq/dogma.HandlesCommandRoute",
			"github.com/dogmatiq/dogma.HandlesEventRoute",
			"github.com/dogmatiq/dogma.ExecutesCommandRoute",
			"github.com/dogmatiq/dogma.SchedulesTimeoutRoute",
			"github.com/dogmatiq/dogma.RecordsEventRoute":

			// At this point we should expect that the interfaces above
			// are produced as a result of calls to following functions:
			// (At the time of writing this code, there is no other way
			// to produce these interfaces)
			//  `github.com/dogmatiq/dogma.HandlesCommand()
			//  `github.com/dogmatiq/dogma.HandlesEvent()`
			//  `github.com/dogmatiq/dogma.ExecutesCommand()`
			//  `github.com/dogmatiq/dogma.RecordsEvent()`
			//  `github.com/dogmatiq/dogma.SchedulesTimeout()`
			if f, ok := mi.X.(*ssa.Call).Common().Value.(*ssa.Function); ok {
				messages.Update(
					message.NameFromStaticType(f.TypeArgs()[0]),
					func(n message.Name, em *configkit.EntityMessage) {
						switch {
						case strings.HasPrefix(f.Name(), "HandlesCommand["):
							em.Kind = message.CommandKind
							em.IsConsumed = true

						case strings.HasPrefix(f.Name(), "HandlesEvent["):
							em.Kind = message.EventKind
							em.IsConsumed = true

						case strings.HasPrefix(f.Name(), "ExecutesCommand["):
							em.Kind = message.CommandKind
							em.IsProduced = true

						case strings.HasPrefix(f.Name(), "RecordsEvent["):
							em.Kind = message.EventKind
							em.IsProduced = true

						case strings.HasPrefix(f.Name(), "SchedulesTimeout["):
							em.Kind = message.TimeoutKind
							em.IsProduced = true
							em.IsConsumed = true
						}
					},
				)

			}
		}
	}
}

// recurseSSAValues recursively traverses the SSA values reachable from val.
// It calls f for each value it encounters. If f returns true, recursion stops.
func recurseSSAValues(
	val ssa.Value,
	seen *[]ssa.Value,
	f func(ssa.Value) bool,
) {
	if val == nil {
		return
	}

	for _, v := range *seen {
		if v == val {
			return
		}
	}

	*seen = append(*seen, val)

	if f(val) {
		return
	}

	// If the value is an instruction, recurse into its operands.
	if instr, ok := val.(ssa.Instruction); ok {
		for _, v := range instr.Operands(nil) {
			recurseSSAValues(*v, seen, f)
		}
	}

	// Check if val has any referrers.
	if referrers := val.Referrers(); referrers != nil {
		for _, instr := range *referrers {

			// If the referrer is a value, recurse into it.
			if v, ok := instr.(ssa.Value); ok {
				recurseSSAValues(v, seen, f)
				continue
			}

			// Otherwise, recurse into the referrer's operands.
			for _, v := range instr.Operands(nil) {
				recurseSSAValues(*v, seen, f)
			}
		}
	}
}

// pkgOfNamedType returns the package in which typ is declared.
//
// It panics if typ is not a named type or a pointer to a named type.
func pkgOfNamedType(typ types.Type) *types.Package {
	pkg, ok := tryPkgOfNamedType(typ)
	if !ok {
		panic(fmt.Sprintf("cannot determine package for anonymous or built-in type %v", typ))
	}

	return pkg
}

// pkgOfNamedType returns the package in which typ is declared.
//
// ok is false if typ is not a named type or a pointer to a named type.
func tryPkgOfNamedType(typ types.Type) (_ *types.Package, ok bool) {
	switch t := typ.(type) {
	case *types.Named:
		return t.Obj().Pkg(), true
	case *types.Alias:
		return t.Obj().Pkg(), true
	case *types.Pointer:
		return tryPkgOfNamedType(t.Elem())
	default:
		return nil, false
	}
}

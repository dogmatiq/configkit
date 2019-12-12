package configkit

import (
	"context"
	"reflect"
	"sync"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// Application is an interface that represents the configuration of a Dogma
// application.
type Application interface {
	Entity

	// Handlers returns the handlers within this application.
	Handlers() HandlerSet
}

// RichApplication is a specialization of Application that exposes information
// about the Go types used to implement the Dogma application.
type RichApplication interface {
	RichEntity

	// Handlers returns the handlers within this application.
	Handlers() HandlerSet

	// RichHandlers returns the handlers within this application.
	RichHandlers() RichHandlerSet

	// Application returns the underlying application.
	Application() dogma.Application
}

// FromApplication returns the configuration for an application.
//
// It panics if the application is configured incorrectly. Use Recover() to
// convert configuration related panic values to errors.
func FromApplication(a dogma.Application) RichApplication {
	cfg := &application{
		entity: entity{
			rt: reflect.TypeOf(a),
		},
		impl: a,
	}

	c := &applicationConfigurer{
		entityConfigurer: entityConfigurer{
			entity: &cfg.entity,
		},
		app: cfg,
	}

	a.Configure(c)

	c.validate()

	return cfg
}

// IsApplicationEqual compares two applications for equality.
//
// It returns true if both applications:
//
//  1. have the same identity
//  2. produce and consume the same messages, with the same roles
//  3. are implemented using the same Go types
//  4. contain equivalent handlers
//
// Point 3. refers to the type used to implement the dogma.Application interface
// (not the type used to implement the configkit.Application interface).
func IsApplicationEqual(a, b Application) bool {
	return a.Identity() == b.Identity() &&
		a.TypeName() == b.TypeName() &&
		a.MessageNames().IsEqual(b.MessageNames()) &&
		a.Handlers().IsEqual(b.Handlers())
}

// ForeignMessageNames returns the subset of message names used by an
// application that must be communicated beyond the scope of that application.
//
// This includes:
//  - commands that are produced by the application, but consumed elsewhere
//  - commands that are consumed by the application, but produced elsewhere
//  - events that are consumed by the application, but produced elsewhere
func ForeignMessageNames(app Application) EntityMessageNames {
	m := app.MessageNames()
	f := EntityMessageNames{
		Produced: message.NameRoles{},
		Consumed: message.NameRoles{},
	}

	for n, r := range m.Produced {
		// Commands MUST always have a handler. Therefore, any command that is
		// produced by this application, but not consumed by this application is
		// considered foreign.
		if r == message.CommandRole && !m.Consumed.Has(n) {
			f.Produced.Add(n, r)
		}
	}

	for n, r := range m.Consumed {
		// Any message, of any role, that is consumed by this application but
		// not produced by this application is considered foreign.
		if !m.Produced.Has(n) {
			f.Consumed.Add(n, r)
		}
	}

	return f
}

// ForeignMessageTypes returns the subset of message types used by an
// application that must be communicated beyond the scope of that application.
//
// This includes:
//	- commands that are produced by this application, but consumed elsewhere
//	- commands that are consumed by this application, but produced elsewhere
//	- events that are consumed by this application, but produced elsewhere
func ForeignMessageTypes(app RichApplication) EntityMessageTypes {
	m := app.MessageTypes()
	f := EntityMessageTypes{
		Produced: message.TypeRoles{},
		Consumed: message.TypeRoles{},
	}

	for t, r := range m.Produced {
		// Commands MUST always have a handler. Therefore, any command that is
		// produced by this application, but not consumed by this application is
		// considered foreign.
		if r == message.CommandRole && !m.Consumed.Has(t) {
			f.Produced.Add(t, r)
		}
	}

	for t, r := range m.Consumed {
		// Any message, of any role, that is consumed by this application but
		// not produced by this application is considered foreign.
		if !m.Produced.Has(t) {
			f.Consumed.Add(t, r)
		}
	}

	return f
}

// application is an implementation of RichApplication.
type application struct {
	entity

	handlers     HandlerSet
	richHandlers RichHandlerSet
	foreignNames EntityMessageNames
	foreignTypes EntityMessageTypes
	impl         dogma.Application
	once         sync.Once
}

func (a *application) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitApplication(ctx, a)
}

func (a *application) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichApplication(ctx, a)
}

func (a *application) Handlers() HandlerSet {
	return a.handlers
}

func (a *application) RichHandlers() RichHandlerSet {
	return a.richHandlers
}

func (a *application) Application() dogma.Application {
	return a.impl
}

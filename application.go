package configkit

import (
	"context"
	"reflect"

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
	cfg, c := fromApplication(a)
	c.mustValidate()
	return cfg
}

func fromApplication(a dogma.Application) (*richApplication, *applicationConfigurer) {
	cfg := &richApplication{
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

	return cfg, c
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

// richApplication is the default implementation of [RichApplication].
type richApplication struct {
	entity

	handlers     HandlerSet
	richHandlers RichHandlerSet
	impl         dogma.Application
}

func (a *richApplication) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitApplication(ctx, a)
}

func (a *richApplication) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichApplication(ctx, a)
}

func (a *richApplication) Handlers() HandlerSet {
	return a.handlers
}

func (a *richApplication) RichHandlers() RichHandlerSet {
	return a.richHandlers
}

func (a *richApplication) Application() dogma.Application {
	return a.impl
}

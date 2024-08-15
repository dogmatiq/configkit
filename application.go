package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
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
	cfg := &richApplication{app: a}
	a.Configure(&applicationConfigurer{config: cfg})

	mustHaveValidIdentity(
		cfg.Identity(),
		cfg.ReflectType(),
	)

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

// richApplication is the default implementation of [RichApplication].
type richApplication struct {
	ident    Identity
	types    EntityMessageTypes
	handlers RichHandlerSet
	app      dogma.Application
}

func (a *richApplication) Identity() Identity {
	return a.ident
}

func (a *richApplication) MessageNames() EntityMessageNames {
	return a.types.asNames()
}

func (a *richApplication) MessageTypes() EntityMessageTypes {
	return a.types
}

func (a *richApplication) TypeName() string {
	return goreflect.NameOf(a.ReflectType())
}

func (a *richApplication) ReflectType() reflect.Type {
	return reflect.TypeOf(a.app)
}

func (a *richApplication) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitApplication(ctx, a)
}

func (a *richApplication) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichApplication(ctx, a)
}

func (a *richApplication) Handlers() HandlerSet {
	return a.handlers.asHandlerSet()
}

func (a *richApplication) RichHandlers() RichHandlerSet {
	return a.handlers
}

func (a *richApplication) Application() dogma.Application {
	return a.app
}

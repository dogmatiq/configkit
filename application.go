package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// Application is an interface that represents the configuration of a Dogma
// application.
type Application interface {
	Entity

	// Handlers returns the handlers within this application.
	Handlers() HandlerSet

	// ForeignMessageNames returns the message names that this application
	// uses that must be communicated to some other Dogma application.
	ForeignMessageNames() message.NameRoles
}

// RichApplication is a specialization of Application that exposes information
// about the Go types used to implement the Dogma application.
type RichApplication interface {
	RichEntity

	// Handlers returns the handlers within this application.
	Handlers() HandlerSet

	// RichHandlers returns the handlers within this application.
	RichHandlers() RichHandlerSet

	// ForeignMessageNames returns the message names that this application
	// uses that must be communicated to some other Dogma application.
	ForeignMessageNames() message.NameRoles

	// ForeignMessageTypes returns the message types that this application
	// uses that must be communicated to some other Dogma application.
	ForeignMessageTypes() message.TypeRoles

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

// application is an implementation of RichApplication.
type application struct {
	entity

	handlers     HandlerSet
	richHandlers RichHandlerSet
	foreignNames message.NameRoles
	foreignTypes message.TypeRoles
	impl         dogma.Application
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

func (a *application) ForeignMessageNames() message.NameRoles {
	return a.foreignNames
}

func (a *application) ForeignMessageTypes() message.TypeRoles {
	return a.foreignTypes
}

func (a *application) Application() dogma.Application {
	return a.impl
}

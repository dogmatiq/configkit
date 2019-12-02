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

	// ForeignMessageNames returns the message names that this application
	// uses that must be communicated beyond the scope of the application.
	//
	// This includes:
	//	- commands that are produced by this application, but consumed elsewhere
	//	- commands that are consumed by this application, but produced elsewhere
	//	- events that are consumed by this application, but produced elsewhere
	ForeignMessageNames() EntityMessageNames
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
	// uses that must be communicated beyond the scope of the application.
	//
	// This includes:
	//	- commands that are produced by this application, but consumed elsewhere
	//	- commands that are consumed by this application, but produced elsewhere
	//	- events that are consumed by this application, but produced elsewhere
	ForeignMessageNames() EntityMessageNames

	// ForeignMessageTypes returns the message types that this application
	// uses that must be communicated beyond the scope of the application.
	//
	// This includes:
	//	- commands that are produced by this application, but consumed elsewhere
	//	- commands that are consumed by this application, but produced elsewhere
	//	- events that are consumed by this application, but produced elsewhere
	ForeignMessageTypes() EntityMessageTypes

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

func (a *application) ForeignMessageNames() EntityMessageNames {
	a.once.Do(a.initForeign)
	return a.foreignNames
}

func (a *application) ForeignMessageTypes() EntityMessageTypes {
	a.once.Do(a.initForeign)
	return a.foreignTypes
}

func (a *application) Application() dogma.Application {
	return a.impl
}

// initForeign initializes a.foreignNames and a.foreignTypes.
func (a *application) initForeign() {
	a.foreignNames = EntityMessageNames{
		Roles:    message.NameRoles{},
		Produced: message.NameRoles{},
		Consumed: message.NameRoles{},
	}

	a.foreignTypes = EntityMessageTypes{
		Roles:    message.TypeRoles{},
		Produced: message.TypeRoles{},
		Consumed: message.TypeRoles{},
	}

	for mt, r := range a.entity.types.Produced {
		if a.entity.types.Consumed.Has(mt) {
			continue
		}

		// Commands MUST always have a handler. Therefore, any command that is
		// produced by this application, but not consumed by this application is
		// considered foreign.
		if r == message.CommandRole {
			a.foreignNames.Roles.Add(mt.Name(), r)
			a.foreignTypes.Roles.Add(mt, r)
			a.foreignNames.Produced.Add(mt.Name(), r)
			a.foreignTypes.Produced.Add(mt, r)
		}
	}

	for mt, r := range a.entity.types.Consumed {
		if a.entity.types.Produced.Has(mt) {
			continue
		}

		// Any message type is considered foreign if it needs to be obtained from
		// elsewhere.
		a.foreignNames.Roles.Add(mt.Name(), r)
		a.foreignTypes.Roles.Add(mt, r)
		a.foreignNames.Consumed.Add(mt.Name(), r)
		a.foreignTypes.Consumed.Add(mt, r)
	}
}

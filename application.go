package configkit

import "github.com/dogmatiq/dogma"

type application interface {
	// Handlers returns the handlers within this application.
	Handlers() []Handler

	// HandlerByIdentity by identity returns the handler with the given identity.
	HandlerByIdentity(Identity) (Handler, bool)

	// HandlerByName by name returns the handler with the given name.
	HandlerByName(string) (Handler, bool)

	// HandlerByKey by name returns the handler with the given key.
	HandlerByKey(string) (Handler, bool)

	// ForeignCommands returns the set of command message types that this
	// application produces but does not consume.
	ForeignCommands() []TypeName

	// ForeignEvents returns the set of event message types that this
	// application consumes but does not produce.
	ForeignEvents() []TypeName
}

// Application is an interface that represents the configuration of a Dogma
// application.
type Application interface {
	Entity
	application
}

// PortableApplication is a specialization of the Application interface that
// does not require the Go types used to implement the Dogma application.
type PortableApplication interface {
	PortableEntity
	application

	// PortableHandlers returns the handlers within this application.
	PortableHandlers() []PortableHandler

	// PortableHandlerByIdentity by identity returns the handler with the given identity.
	PortableHandlerByIdentity(Identity) (PortableHandler, bool)

	// PortableHandlerByName by name returns the handler with the given name.
	PortableHandlerByName(string) (PortableHandler, bool)

	// PortableHandlerByKey by name returns the handler with the given key.
	PortableHandlerByKey(string) (PortableHandler, bool)
}

// RichApplication is a specialization of the Application interface that has
// access to the Go types used to implement the Dogma application.
type RichApplication interface {
	RichEntity
	application

	// Application returns the application implementation.
	Application() dogma.Application

	// RichHandlers returns the handlers within this application.
	RichHandlers() []RichHandler

	// RichHandlerByIdentity by identity returns the handler with the given identity.
	RichHandlerByIdentity(Identity) (RichHandler, bool)

	// RichHandlerByName by name returns the handler with the given name.
	RichHandlerByName(string) (RichHandler, bool)

	// RichHandlerByKey by name returns the handler with the given key.
	RichHandlerByKey(string) (RichHandler, bool)
}

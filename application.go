package configkit

import "github.com/dogmatiq/dogma"

// Application is an interface that represents the configuration of a Dogma
// application.
type Application interface {
	Entity

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

// RichApplication is a specialization of the Application interface that has
// access to the Go types used to implement the Dogma application.
type RichApplication interface {
	RichEntity

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

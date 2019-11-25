package configkit

import (
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

	// Application returns the underlying application .
	Application() dogma.Application
}

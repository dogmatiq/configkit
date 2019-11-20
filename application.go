package configkit

import "github.com/dogmatiq/dogma"

type application interface {
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
}

// RichApplication is a specialization of the Application interface that has
// access to the Go types used to implement the Dogma application.
type RichApplication interface {
	RichEntity
	application

	// Application returns the application implementation.
	Application() dogma.Application
}

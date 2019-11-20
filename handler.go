package configkit

type handler interface {
	// HandlerType returns the type of handler.
	HandlerType() HandlerType
}

// Handler is an interface that represents the configuration of a Dogma message
// handler.
type Handler interface {
	Entity
	handler
}

// PortableHandler is a specialization of the HandlerConfig interface that
// does not require the Go types used to implement the Dogma application.
type PortableHandler interface {
	PortableEntity
	handler
}

// RichHandler is a specialization of the HandlerConfig interface that has
// access to the Go types used to implement the Dogma application.
type RichHandler interface {
	RichEntity
	handler
}

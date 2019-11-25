package configkit

// Handler is a specialization of the Entity interface for message handlers.
type Handler interface {
	Entity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType
}

// RichHandler is a specialization of the Handler interface that exposes
// information about the Go types used to implement the Dogma application.
type RichHandler interface {
	RichEntity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType
}

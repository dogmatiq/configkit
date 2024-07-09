package configkit

// Handler is a specialization of the Entity interface for message handlers.
type Handler interface {
	Entity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	IsDisabled() bool
}

// RichHandler is a specialization of the Handler interface that exposes
// information about the Go types used to implement the Dogma application.
type RichHandler interface {
	RichEntity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	IsDisabled() bool
}

// IsHandlerEqual compares two handlers for equality.
//
// It returns true if both handlers:
//
//  1. have the same identity
//  2. produce and consume the same messages, with the same roles
//  3. are implemented using the same Go types
//
// Point 3. refers to the type used to implement the dogma.Aggregate,
// dogma.Process, dogma.Integration or dogma.Projection interface (not the type
// used to implement the configkit.Handler interface).
//
// This definition of equality relies on the fact that no single Go type can
// implement more than one Dogma handler interface because they all contain
// a Configure() method with different signatures.
func IsHandlerEqual(a, b Handler) bool {
	return a.Identity() == b.Identity() &&
		a.TypeName() == b.TypeName() &&
		a.HandlerType() == b.HandlerType() &&
		a.MessageNames().IsEqual(b.MessageNames())
}

type handler struct {
	entity
	isDisabled bool
}

func (h *handler) IsDisabled() bool {
	return h.isDisabled
}

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

// RichApplication is an implementation of Application that has access to the Go
// types used to implement the Dogma application.
type RichApplication struct {
	Application dogma.Application
}

// Identity returns the identity of the entity.
func (*RichApplication) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichApplication) TypeName() string {
	panic("not implemented")
}

// MessageNames returns information about the messages used by the entity.
func (*RichApplication) MessageNames() EntityMessageNames {
	panic("not implemented")
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (*RichApplication) AcceptVisitor(ctx context.Context, v Visitor) error {
	panic("not implemented")
}

// ReflectType returns the reflect.Type of the Dogma entity.
func (*RichApplication) ReflectType() reflect.Type {
	panic("not implemented")
}

// MessageTypes returns information about the messages used by the entity.
func (*RichApplication) MessageTypes() EntityMessageTypes {
	panic("not implemented")
}

// AcceptRichVisitor calls the appropriate method on v for this
// configuration type.
func (*RichApplication) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	panic("not implemented")
}

// Handlers returns the handlers within this application.
func (*RichApplication) Handlers() HandlerSet {
	panic("not implemented")
}

// ForeignMessageNames returns the message names that this application
// uses that must be communicated to some other Dogma application.
func (*RichApplication) ForeignMessageNames() message.NameRoles {
	panic("not implemented")
}

// RichHandlers returns the handlers within this application.
func (*RichApplication) RichHandlers() RichHandlerSet {
	panic("not implemented")
}

// ForeignMessageTypes returns the message types that this application
// uses that must be communicated to some other Dogma application.
func (*RichApplication) ForeignMessageTypes() message.TypeRoles {
	panic("not implemented")
}

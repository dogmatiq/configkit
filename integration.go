package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Integration is an interface that represents the configuration of a Dogma
// integration message handler.
type Integration interface {
	Handler
}

// RichIntegration is an implementation of Integration that exposes information
// about the Go types used to implement the underlying Dogma handler.
type RichIntegration struct {
	Handler dogma.IntegrationMessageHandler
}

// Identity returns the identity of the entity.
func (*RichIntegration) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichIntegration) TypeName() string {
	panic("not implemented")
}

// MessageNames returns information about the messages used by the entity.
func (*RichIntegration) MessageNames() EntityMessageNames {
	panic("not implemented")
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (*RichIntegration) AcceptVisitor(ctx context.Context, v Visitor) error {
	panic("not implemented")
}

// ReflectType returns the reflect.Type of the Dogma entity.
func (*RichIntegration) ReflectType() reflect.Type {
	panic("not implemented")
}

// MessageTypes returns information about the messages used by the entity.
func (*RichIntegration) MessageTypes() EntityMessageTypes {
	panic("not implemented")
}

// AcceptRichVisitor calls the appropriate method on v for this
// configuration type.
func (*RichIntegration) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	panic("not implemented")
}

// HandlerType returns the type of handler.
func (*RichIntegration) HandlerType() HandlerType {
	panic("not implemented")
}

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

// RichIntegration is an implementation of Integration that has access to the Go
// types used to implement the underlying Dogma handler.
type RichIntegration struct {
	Handler dogma.IntegrationMessageHandler
}

// Identity returns the identity of the entity.
func (*RichIntegration) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichIntegration) TypeName() TypeName {
	panic("not implemented")
}

// Messages returns the messages used by the entity in any way.
func (*RichIntegration) Messages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ConsumedMessages returns the message types consumed by the entity.
func (*RichIntegration) ConsumedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ProducedMessages returns the message types produced by the entity.
func (*RichIntegration) ProducedMessages() map[TypeName]MessageRole {
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

// RichMessages returns the messages used by the entity in any way.
func (*RichIntegration) RichMessages() map[MessageType]MessageRole {
	panic("not implemented")
}

// RichConsumedMessages returns the message types consumed by the entity.
func (*RichIntegration) RichConsumedMessages() map[MessageType]MessageRole {
	panic("not implemented")
}

// RichProducedMessages returns the message types produced by the entity.
func (*RichIntegration) RichProducedMessages() map[MessageType]MessageRole {
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

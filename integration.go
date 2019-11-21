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

// ReflectTypeOf returns the reflect.Type of the type with the given name.
// It panics if the type name is not used within the entity.
func (*RichIntegration) ReflectTypeOf(TypeName) reflect.Type {
	panic("not implemented")
}

// MessageTypeOf returns the MessageType of the type with the given name.
// It panics if the type is not used as a message within the entity.
func (*RichIntegration) MessageTypeOf(TypeName) (MessageType, bool) {
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

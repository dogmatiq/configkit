package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Process is an interface that represents the configuration of a Dogma process
// message handler.
type Process interface {
	Handler
}

// RichProcess is an implementation of Process that has access to the Go types
// used to implement the underlying Dogma handler.
type RichProcess struct {
	Handler dogma.ProcessMessageHandler
}

// Identity returns the identity of the entity.
func (*RichProcess) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichProcess) TypeName() TypeName {
	panic("not implemented")
}

// ConsumedMessages returns the message types consumed by the entity.
func (*RichProcess) ConsumedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ProducedMessages returns the message types produced by the entity.
func (*RichProcess) ProducedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (*RichProcess) AcceptVisitor(ctx context.Context, v Visitor) error {
	panic("not implemented")
}

// ReflectType returns the reflect.Type of the Dogma entity.
func (*RichProcess) ReflectType() reflect.Type {
	panic("not implemented")
}

// ReflectTypeOf returns the reflect.Type of the type with the given name.
// It panics if the type name is not used within the entity.
func (*RichProcess) ReflectTypeOf(TypeName) reflect.Type {
	panic("not implemented")
}

// MessageTypeOf returns the MessageType of the type with the given name.
// It panics if the type is not used as a message within the entity.
func (*RichProcess) MessageTypeOf(TypeName) (MessageType, bool) {
	panic("not implemented")
}

// AcceptRichVisitor calls the appropriate method on v for this
// configuration type.
func (*RichProcess) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	panic("not implemented")
}

// HandlerType returns the type of handler.
func (*RichProcess) HandlerType() HandlerType {
	panic("not implemented")
}

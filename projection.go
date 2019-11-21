package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Projection is an interface that represents the configuration of a Dogma
// projection message handler.
type Projection interface {
	Handler
}

// RichProjection is an implementation of Projection that has access to the Go
// types used to implement the underlying Dogma handler.
type RichProjection struct {
	Handler dogma.ProjectionMessageHandler
}

// Identity returns the identity of the entity.
func (*RichProjection) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichProjection) TypeName() TypeName {
	panic("not implemented")
}

// Messages returns the messages used by the entity in any way.
func (*RichProjection) Messages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ConsumedMessages returns the message types consumed by the entity.
func (*RichProjection) ConsumedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ProducedMessages returns the message types produced by the entity.
func (*RichProjection) ProducedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (*RichProjection) AcceptVisitor(ctx context.Context, v Visitor) error {
	panic("not implemented")
}

// ReflectType returns the reflect.Type of the Dogma entity.
func (*RichProjection) ReflectType() reflect.Type {
	panic("not implemented")
}

// ReflectTypeOf returns the reflect.Type of the type with the given name.
// It panics if the type name is not used within the entity.
func (*RichProjection) ReflectTypeOf(TypeName) reflect.Type {
	panic("not implemented")
}

// MessageTypeOf returns the MessageType of the type with the given name.
// It panics if the type is not used as a message within the entity.
func (*RichProjection) MessageTypeOf(TypeName) (MessageType, bool) {
	panic("not implemented")
}

// AcceptRichVisitor calls the appropriate method on v for this
// configuration type.
func (*RichProjection) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	panic("not implemented")
}

// HandlerType returns the type of handler.
func (*RichProjection) HandlerType() HandlerType {
	panic("not implemented")
}

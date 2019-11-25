package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Aggregate is an interface that represents the configuration of a Dogma
// aggregate message handler.
type Aggregate interface {
	Handler
}

// RichAggregate is an implementation of Aggregate that has access to the Go
// types used to implement the underlying Dogma handler.
type RichAggregate struct {
	Handler dogma.AggregateMessageHandler
}

// Identity returns the identity of the entity.
func (*RichAggregate) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichAggregate) TypeName() TypeName {
	panic("not implemented")
}

// Messages returns the messages used by the entity in any way.
func (*RichAggregate) Messages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ConsumedMessages returns the message types consumed by the entity.
func (*RichAggregate) ConsumedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ProducedMessages returns the message types produced by the entity.
func (*RichAggregate) ProducedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// AcceptVisitor calls the appropriate method on v for this entity type.
func (*RichAggregate) AcceptVisitor(ctx context.Context, v Visitor) error {
	panic("not implemented")
}

// ReflectType returns the reflect.Type of the Dogma entity.
func (*RichAggregate) ReflectType() reflect.Type {
	panic("not implemented")
}

// RichMessages returns the messages used by the entity in any way.
func (*RichAggregate) RichMessages() map[MessageType]MessageRole {
	panic("not implemented")
}

// RichConsumedMessages returns the message types consumed by the entity.
func (*RichAggregate) RichConsumedMessages() map[MessageType]MessageRole {
	panic("not implemented")
}

// RichProducedMessages returns the message types produced by the entity.
func (*RichAggregate) RichProducedMessages() map[MessageType]MessageRole {
	panic("not implemented")
}

// AcceptRichVisitor calls the appropriate method on v for this
// configuration type.
func (*RichAggregate) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	panic("not implemented")
}

// HandlerType returns the type of handler.
func (*RichAggregate) HandlerType() HandlerType {
	panic("not implemented")
}

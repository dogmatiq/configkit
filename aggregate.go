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

// RichAggregate is an implementation of Aggregate that exposes informatiom
// about the Go types used to implement the underlying Dogma handler.
type RichAggregate struct {
	Handler dogma.AggregateMessageHandler
}

// Identity returns the identity of the entity.
func (*RichAggregate) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichAggregate) TypeName() string {
	panic("not implemented")
}

// MessageNames returns information about the messages used by the entity.
func (*RichAggregate) MessageNames() EntityMessageNames {
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

// MessageTypes returns information about the messages used by the entity.
func (*RichAggregate) MessageTypes() EntityMessageTypes {
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

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

// RichProjection is an implementation of Projection that exposes information
// about the Go types used to implement the underlying Dogma handler.
type RichProjection struct {
	Handler dogma.ProjectionMessageHandler
}

// Identity returns the identity of the entity.
func (*RichProjection) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichProjection) TypeName() string {
	panic("not implemented")
}

// MessageNames returns information about the messages used by the entity.
func (*RichProjection) MessageNames() EntityMessageNames {
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

// MessageTypes returns information about the messages used by the entity.
func (*RichProjection) MessageTypes() EntityMessageTypes {
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

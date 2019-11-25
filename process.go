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

// RichProcess is an implementation of Process that exposes information about
// the Go types used to implement the underlying Dogma handler.
type RichProcess struct {
	Handler dogma.ProcessMessageHandler
}

// Identity returns the identity of the entity.
func (*RichProcess) Identity() Identity {
	panic("not implemented")
}

// TypeName returns the fully-qualified type name of the entity.
func (*RichProcess) TypeName() string {
	panic("not implemented")
}

// MessageNames returns information about the messages used by the entity.
func (*RichProcess) MessageNames() EntityMessageNames {
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

// MessageTypes returns information about the messages used by the entity.
func (*RichProcess) MessageTypes() EntityMessageTypes {
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

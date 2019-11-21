package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Application is an interface that represents the configuration of a Dogma
// application.
type Application interface {
	Entity

	// Handlers returns the handlers within this application.
	Handlers() []Handler

	// HandlerByIdentity by identity returns the handler with the given identity.
	HandlerByIdentity(Identity) (Handler, bool)

	// HandlerByName by name returns the handler with the given name.
	HandlerByName(string) (Handler, bool)

	// HandlerByKey by name returns the handler with the given key.
	HandlerByKey(string) (Handler, bool)

	// ConsumersOf returns the handlers that consume messages of the given type.
	ConsumersOf(TypeName) []Handler

	// ProducersOf returns the handlers that produce messages of the given type.
	ProducersOf(TypeName) []Handler

	// ForeignMessages returns the set of message types that this application
	// references that originate or are destined for some other application.
	ForeignMessages() map[TypeName]MessageRole
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
func (*RichApplication) TypeName() TypeName {
	panic("not implemented")
}

// Messages returns the messages used by the entity in any way.
func (*RichApplication) Messages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ConsumedMessages returns the message types consumed by the entity.
func (*RichApplication) ConsumedMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// ProducedMessages returns the message types produced by the entity.
func (*RichApplication) ProducedMessages() map[TypeName]MessageRole {
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

// ReflectTypeOf returns the reflect.Type of the type with the given name.
// It panics if the type name is not used within the entity.
func (*RichApplication) ReflectTypeOf(TypeName) reflect.Type {
	panic("not implemented")
}

// MessageTypeOf returns the MessageType of the type with the given name.
// It panics if the type is not used as a message within the entity.
func (*RichApplication) MessageTypeOf(TypeName) (MessageType, bool) {
	panic("not implemented")
}

// AcceptRichVisitor calls the appropriate method on v for this
// configuration type.
func (*RichApplication) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	panic("not implemented")
}

// Handlers returns the handlers within this application.
func (*RichApplication) Handlers() []Handler {
	panic("not implemented")
}

// HandlerByIdentity by identity returns the handler with the given identity.
func (*RichApplication) HandlerByIdentity(Identity) (Handler, bool) {
	panic("not implemented")
}

// HandlerByName by name returns the handler with the given name.
func (*RichApplication) HandlerByName(string) (Handler, bool) {
	panic("not implemented")
}

// HandlerByKey by name returns the handler with the given key.
func (*RichApplication) HandlerByKey(string) (Handler, bool) {
	panic("not implemented")
}

// ConsumersOf returns the handlers that consume messages of the given type.
func (*RichApplication) ConsumersOf(TypeName) []Handler {
	panic("not implemented")
}

// ProducersOf returns the handlers that produce messages of the given type.
func (*RichApplication) ProducersOf(TypeName) []Handler {
	panic("not implemented")
}

// ForeignMessages returns the set of message types that this application
// references that originate or are destined for some other application.
func (*RichApplication) ForeignMessages() map[TypeName]MessageRole {
	panic("not implemented")
}

// RichHandlers returns the handlers within this application.
func (*RichApplication) RichHandlers() []RichHandler {
	panic("not implemented")
}

// RichHandlerByIdentity by identity returns the handler with the given identity.
func (*RichApplication) RichHandlerByIdentity(Identity) (RichHandler, bool) {
	panic("not implemented")
}

// RichHandlerByName by name returns the handler with the given name.
func (*RichApplication) RichHandlerByName(string) (RichHandler, bool) {
	panic("not implemented")
}

// RichHandlerByKey by name returns the handler with the given key.
func (*RichApplication) RichHandlerByKey(string) (RichHandler, bool) {
	panic("not implemented")
}

// RichConsumersOf returns the handlers that consume messages of the given type.
func (*RichApplication) RichConsumersOf(TypeName) []*RichHandler {
	panic("not implemented")
}

// RichProducersOf returns the handlers that produce messages of the given type.
func (*RichApplication) RichProducersOf(TypeName) []*RichHandler {
	panic("not implemented")
}

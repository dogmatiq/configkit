package configkit

import (
	"context"
	"reflect"
)

// TypeName is a fully-qualified name for a Go type.
type TypeName string

// Entity is an interface that represents the configuration of a Dogma "entity"
// such as an application or handler.
//
// Each implementation of this interface represents the configuration described
// by a call to the entity's Configure() method.
type Entity interface {
	// Identity returns the identity of the entity.
	Identity() Identity

	// TypeName returns the fully-qualified type name of the entity.
	TypeName() TypeName

	// ConsumedMessages returns the message types consumed by the entity.
	ConsumedMessages() map[TypeName]MessageRole

	// ProducedMessages returns the message types produced by the entity.
	ProducedMessages() map[TypeName]MessageRole

	// AcceptVisitor calls the appropriate method on v for this entity type.
	AcceptVisitor(ctx context.Context, v Visitor) error
}

// PortableEntity is a specialization of the Entity interface that does not
// require the Go types used to implement the Dogma entity.
type PortableEntity interface {
	Entity

	// AcceptPortableVisitor calls the appropriate method on v for this
	// configuration type.
	AcceptPortableVisitor(ctx context.Context, v PortableVisitor) error
}

// RichEntity is a specialization of the Entity interface that has access to the
// Go types used to implement the Dogma entity.
type RichEntity interface {
	Entity

	// ReflectType returns the reflect.Type of the Dogma entity.
	ReflectType() reflect.Type

	// ReflectTypeOf returns the reflect.Type of the type with the given name.
	// It panics if the type name is not used within the entity.
	ReflectTypeOf(TypeName) reflect.Type

	// MessageTypeOf returns the MessageType of the type with the given name.
	// It panics if the type is not used as a message within the entity.
	MessageTypeOf(TypeName) (MessageType, bool)

	// AcceptRichVisitor calls the appropriate method on v for this
	// configuration type.
	AcceptRichVisitor(ctx context.Context, v RichVisitor) error
}

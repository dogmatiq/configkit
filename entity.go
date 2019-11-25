package configkit

import (
	"context"
	"reflect"
)

// Entity is an interface that represents the configuration of a Dogma "entity"
// such as an application or handler.
//
// Each implementation of this interface represents the configuration described
// by a call to the entity's Configure() method.
type Entity interface {
	// Identity returns the identity of the entity.
	Identity() Identity

	// TypeName returns the fully-qualified type name of the entity.
	TypeName() string

	// MessageNames returns information about the messages used by the entity.
	MessagesNames() MessageNames

	// AcceptVisitor calls the appropriate method on v for this entity type.
	AcceptVisitor(ctx context.Context, v Visitor) error
}

// RichEntity is a specialization of the Entity interface that exposes
// information about the Go types used to implement the Dogma entity.
type RichEntity interface {
	Entity

	// ReflectType returns the reflect.Type of the Dogma entity.
	ReflectType() reflect.Type

	// MessageTypes returns information about the messages used by the entity.
	MessageTypes() MessageTypes

	// AcceptRichVisitor calls the appropriate method on v for this
	// configuration type.
	AcceptRichVisitor(ctx context.Context, v RichVisitor) error
}

package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename"
	"github.com/dogmatiq/configkit/message"
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
	MessageNames() EntityMessageNames

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
	MessageTypes() EntityMessageTypes

	// AcceptRichVisitor calls the appropriate method on v for this
	// configuration type.
	AcceptRichVisitor(ctx context.Context, v RichVisitor) error
}

// EntityMessageNames describes how messages are used within a Dogma entity
// where each message is identified by its name.
type EntityMessageNames struct {
	// Roles is a map of message name to its role within the entity.
	Roles message.NameRoles

	// Produced is a set of message names produced by the entity.
	Produced message.NameRoles

	// Consumed is a set of message names consumed by the entity.
	Consumed message.NameRoles
}

// EntityMessageTypes describes how messages are used within a Dogma entity
// where each message is identified by its type.
type EntityMessageTypes struct {
	// Roles is a map of message type to its role within the entity.
	Roles message.TypeRoles

	// Produced is a set of message types produced by the entity.
	Produced message.TypeRoles

	// Consumed is a set of message types consumed by the entity.
	Consumed message.TypeRoles
}

// IsEqual returns true if m is equal to o.
func (m EntityMessageTypes) IsEqual(o EntityMessageTypes) bool {
	return m.Roles.IsEqual(o.Roles) &&
		m.Produced.IsEqual(o.Produced) &&
		m.Consumed.IsEqual(o.Consumed)
}

// entity is a partial implementation of RichEntity.
type entity struct {
	rt reflect.Type

	ident Identity
	names EntityMessageNames
	types EntityMessageTypes
}

func (e *entity) Identity() Identity {
	return e.ident
}

func (e *entity) MessageNames() EntityMessageNames {
	return e.names
}

func (e *entity) MessageTypes() EntityMessageTypes {
	return e.types
}

func (e *entity) TypeName() string {
	return typename.Of(e.ReflectType())
}

func (e *entity) ReflectType() reflect.Type {
	return e.rt
}

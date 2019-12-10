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

// IsEqual returns true if m is equal to o.
func (m EntityMessageNames) IsEqual(o EntityMessageNames) bool {
	return m.Roles.IsEqual(o.Roles) &&
		m.Produced.IsEqual(o.Produced) &&
		m.Consumed.IsEqual(o.Consumed)
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

// IsEqual compares two entities for equality.
//
// It returns true if both entities:
//
//  1. have the same identity
//  2. produce and consume the same messages, with the same roles
//  3. are implemented using the same Go types
//
// Point 3. refers to the type used to implement the dogma.Application,
// dogma.Aggregate, dogma.Process, dogma.Integration or dogma.Projection
// interface (not the type used to implement the Entity interface).
//
// This definition of equality relies on the fact that no single Go type can
// implement more than one these interfaces because they all contain a
// Configure() method with different signatures.
func IsEqual(a, b Entity) bool {
	if a.Identity() != b.Identity() {
		return false
	}

	if a.TypeName() != b.TypeName() {
		return false
	}

	return a.MessageNames().IsEqual(b.MessageNames())
}

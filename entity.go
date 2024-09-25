package configkit

import (
	"context"
	"reflect"

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
	// Produced is a set of message names produced by the entity.
	Produced message.NameRoles

	// Consumed is a set of message names consumed by the entity.
	Consumed message.NameRoles
}

// RoleOf returns the role associated with n, if any.
func (m EntityMessageNames) RoleOf(n message.Name) (message.Role, bool) {
	if r, ok := m.Produced[n]; ok {
		return r, true
	}

	r, ok := m.Consumed[n]
	return r, ok
}

// All returns the type roles of all messages, both produced and consumed.
func (m EntityMessageNames) All() message.NameRoles {
	roles := message.NameRoles{}

	for n, r := range m.Produced {
		roles[n] = r
	}

	for n, r := range m.Consumed {
		roles[n] = r
	}

	return roles
}

// IsEqual returns true if m is equal to o.
func (m EntityMessageNames) IsEqual(o EntityMessageNames) bool {
	return m.Produced.IsEqual(o.Produced) &&
		m.Consumed.IsEqual(o.Consumed)
}

// EntityMessageTypes describes how messages are used within a Dogma entity
// where each message is identified by its type.
type EntityMessageTypes struct {
	// Produced is a set of message types produced by the entity.
	Produced message.TypeRoles

	// Consumed is a set of message types consumed by the entity.
	Consumed message.TypeRoles
}

// RoleOf returns the role associated with t, if any.
func (m EntityMessageTypes) RoleOf(t message.Type) (message.Role, bool) {
	if r, ok := m.Produced[t]; ok {
		return r, true
	}

	r, ok := m.Consumed[t]
	return r, ok
}

// All returns the type roles of all messages, both produced and consumed.
func (m EntityMessageTypes) All() message.TypeRoles {
	roles := message.TypeRoles{}

	for t, r := range m.Produced {
		roles[t] = r
	}

	for t, r := range m.Consumed {
		roles[t] = r
	}

	return roles
}

// IsEqual returns true if m is equal to o.
func (m EntityMessageTypes) IsEqual(o EntityMessageTypes) bool {
	return m.Produced.IsEqual(o.Produced) &&
		m.Consumed.IsEqual(o.Consumed)
}

func (m EntityMessageTypes) asNames() EntityMessageNames {
	var names EntityMessageNames

	if len(m.Produced) != 0 {
		names.Produced = make(message.NameRoles, len(m.Produced))
		for t, r := range m.Produced {
			names.Produced.Add(t.Name(), r)
		}
	}

	if len(m.Consumed) != 0 {
		names.Consumed = make(message.NameRoles, len(m.Consumed))
		for t, r := range m.Consumed {
			names.Consumed.Add(t.Name(), r)
		}
	}

	return names
}

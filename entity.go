package configkit

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/collection/sets"
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

// EntityMessageNames describes the messages used by a Dogma entity where each
// message is identified by its name.
type EntityMessageNames struct {
	// Kinds is a map of message type name to that type's kind.
	Kinds map[message.Name]message.Kind

	// Produced contains the names of the messages produced by the entity.
	Produced sets.Set[message.Name]

	// Consumed contains the names of the messages consumed by the entity.
	Consumed sets.Set[message.Name]
}

// Has returns true if entity uses a message type with the given name.
func (names EntityMessageNames) Has(n message.Name) bool {
	return names.Produced.Has(n) || names.Consumed.Has(n)
}

// IsEqual returns true if names is equal to n.
func (names EntityMessageNames) IsEqual(n EntityMessageNames) bool {
	if len(names.Kinds) != len(n.Kinds) {
		return false
	}

	for name, k := range names.Kinds {
		if x, ok := n.Kinds[name]; !ok || x != k {
			return false
		}
	}

	return names.Produced.IsEqual(n.Produced) &&
		names.Consumed.IsEqual(n.Consumed)
}

func (names *EntityMessageNames) union(n EntityMessageNames) {
	merge := func(dst, src *sets.Set[message.Name]) {
		for name := range src.All() {
			k, ok := n.Kinds[name]
			if !ok {
				panic(fmt.Sprintf(
					"message type %q has no associated message kind",
					name,
				))
			}

			if x, ok := names.Kinds[name]; !ok {
				if names.Kinds == nil {
					names.Kinds = map[message.Name]message.Kind{}
				}
				names.Kinds[name] = k
			} else if x != k {
				panic(fmt.Sprintf(
					"message type %q has conflicting kinds %s and %s",
					name,
					x,
					k,
				))
			}

			dst.Add(name)
		}
	}

	merge(&names.Produced, &n.Produced)
	merge(&names.Consumed, &n.Consumed)
}

// EntityMessageTypes describes the message types used by a Dogma entity.
type EntityMessageTypes struct {
	// Produced is a set of message types produced by the entity.
	Produced sets.Set[message.Type]

	// Consumed is a set of message types consumed by the entity.
	Consumed sets.Set[message.Type]
}

// Has returns true if the entity uses messages of the given type.
func (types EntityMessageTypes) Has(t message.Type) bool {
	return types.Produced.Has(t) || types.Consumed.Has(t)
}

// IsEqual returns true if types is equal to t.
func (types EntityMessageTypes) IsEqual(t EntityMessageTypes) bool {
	return types.Produced.IsEqual(t.Produced) &&
		types.Consumed.IsEqual(t.Consumed)
}

func (types *EntityMessageTypes) union(t EntityMessageTypes) {
	for mt := range t.Produced.All() {
		types.Produced.Add(mt)
	}

	for mt := range t.Consumed.All() {
		types.Consumed.Add(mt)
	}
}

func (types EntityMessageTypes) asNames() EntityMessageNames {
	var names EntityMessageNames

	for t := range types.Produced.All() {
		if names.Kinds == nil {
			names.Kinds = map[message.Name]message.Kind{}
		}

		n := t.Name()
		names.Kinds[n] = t.Kind()
		names.Produced.Add(n)
	}

	for t := range types.Consumed.All() {
		if names.Kinds == nil {
			names.Kinds = map[message.Name]message.Kind{}
		}

		n := t.Name()
		names.Kinds[n] = t.Kind()
		names.Consumed.Add(n)
	}

	return names
}

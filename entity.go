package configkit

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"slices"

	"github.com/dogmatiq/enginekit/message"
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
	MessageNames() EntityMessages[message.Name]

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
	MessageTypes() EntityMessages[message.Type]

	// AcceptRichVisitor calls the appropriate method on v for this
	// configuration type.
	AcceptRichVisitor(ctx context.Context, v RichVisitor) error
}

// EntityMessage describes a message used by a Dogma entity.
type EntityMessage struct {
	Kind                   message.Kind
	IsProduced, IsConsumed bool
}

// EntityMessages describes the messages used by a Dogma entity.
type EntityMessages[K comparable] map[K]EntityMessage

// IsEqual returns true if m is equal to n.
func (m EntityMessages[K]) IsEqual(n EntityMessages[K]) bool {
	if len(m) != len(n) {
		return false
	}

	for k, v := range m {
		if x, ok := n[k]; !ok || x != v {
			return false
		}
	}

	return true
}

// Produced returns an iterator that yields the messages that are produced by
// the entity.
func (m EntityMessages[K]) Produced(filter ...message.Kind) iter.Seq2[K, message.Kind] {
	return func(yield func(K, message.Kind) bool) {
		for k, v := range m {
			if v.IsProduced {
				if len(filter) == 0 || slices.Contains(filter, v.Kind) {
					if !yield(k, v.Kind) {
						return
					}
				}
			}
		}
	}
}

// Consumed returns an iterator that yields the messages that are consumed by
// the entity.
func (m EntityMessages[K]) Consumed(filter ...message.Kind) iter.Seq2[K, message.Kind] {
	return func(yield func(K, message.Kind) bool) {
		for n, m := range m {
			if m.IsConsumed {
				if len(filter) == 0 || slices.Contains(filter, m.Kind) {
					if !yield(n, m.Kind) {
						return
					}
				}
			}
		}
	}
}

// Update updates the message with the given key by calling fn.
//
// If, after calling fn, the [EntityMessage] is neither produced nor consumed,
// it is removed from the map.
func (m EntityMessages[K]) Update(k K, fn func(K, *EntityMessage)) {
	em, ok := m[k]

	fn(k, &em)

	if em.IsConsumed || em.IsProduced {
		m[k] = em
	} else if ok {
		delete(m, k)
	}
}

func (m EntityMessages[K]) merge(n EntityMessages[K]) {
	for k, em := range n {
		x, ok := m[k]

		if !ok {
			x.Kind = em.Kind
		} else if x.Kind != em.Kind {
			panic(fmt.Sprintf("message %v has conflicting kinds %s and %s", k, x.Kind, em.Kind))
		}

		if em.IsProduced {
			x.IsProduced = true
		}

		if em.IsConsumed {
			x.IsConsumed = true
		}

		m[k] = x
	}
}

func asMessageNames(types EntityMessages[message.Type]) EntityMessages[message.Name] {
	names := make(EntityMessages[message.Name], len(types))

	for t, em := range types {
		names[t.Name()] = em
	}

	return names
}

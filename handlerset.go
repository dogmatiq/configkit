package configkit

import (
	"context"

	"github.com/dogmatiq/configkit/message"
)

// HandlerSet is a collection of handlers.
type HandlerSet map[Identity]Handler

// NewHandlerSet returns a HandlerSet containing the given handlers.
//
// It panics if any of the handler identities conflict.
func NewHandlerSet(handlers ...Handler) HandlerSet {
	s := HandlerSet{}

	for _, h := range handlers {
		if !s.Add(h) {
			panic("handler set contains conflicting identities")
		}
	}

	return s
}

// Add adds a handler to the set.
//
// It returns true if the handler was added, or false if the set already
// contained a handler with the same name or key as h.
func (s HandlerSet) Add(h Handler) bool {
	i := h.Identity()
	for x := range s {
		if i.ConflictsWith(x) {
			return false
		}
	}

	s[i] = h
	return true
}

// Has returns true if s contains h.
func (s HandlerSet) Has(h Handler) bool {
	x, ok := s[h.Identity()]
	return ok && x == h
}

// ByIdentity returns the handler with the given identity.
func (s HandlerSet) ByIdentity(i Identity) (Handler, bool) {
	h, ok := s[i]
	return h, ok
}

// ByName returns the handler with the given name.
func (s HandlerSet) ByName(n string) (Handler, bool) {
	return s.Find(func(h Handler) bool {
		return h.Identity().Name == n
	})
}

// ByKey returns the handler with the given key.
func (s HandlerSet) ByKey(k string) (Handler, bool) {
	return s.Find(func(h Handler) bool {
		return h.Identity().Key == k
	})
}

// ByType returns the subset of handlers of the given type.
func (s HandlerSet) ByType(t HandlerType) HandlerSet {
	return s.Filter(func(h Handler) bool {
		return h.HandlerType() == t
	})
}

// ConsumersOf returns the subset of handlers that consume messages with the
// given name.
func (s HandlerSet) ConsumersOf(n message.Name) HandlerSet {
	return s.Filter(func(h Handler) bool {
		return h.MessageNames().Consumed.Has(n)
	})
}

// ProducersOf returns the subset of handlers that produce messages with the
// given name.
func (s HandlerSet) ProducersOf(n message.Name) HandlerSet {
	return s.Filter(func(h Handler) bool {
		return h.MessageNames().Produced.Has(n)
	})
}

// IsEqual returns true if o contains the same handlers as s.
func (s HandlerSet) IsEqual(o HandlerSet) bool {
	if len(s) != len(o) {
		return false
	}

	for i, h := range s {
		x, ok := o[i]
		if !ok || !IsEqual(x, h) {
			return false
		}
	}

	return true
}

// Find returns a handler from the set for which the given predicate function
// returns true.
func (s HandlerSet) Find(fn func(Handler) bool) (Handler, bool) {
	for _, h := range s {
		if fn(h) {
			return h, true
		}
	}

	return nil, false
}

// Filter returns the subset of handlers for which the given predicate function
// returns true.
func (s HandlerSet) Filter(fn func(Handler) bool) HandlerSet {
	subset := HandlerSet{}

	for i, h := range s {
		if fn(h) {
			subset[i] = h
		}
	}

	return subset
}

// AcceptVisitor visits each handler in the set.
//
// It returns the error returned by the first handler to return a non-nil error.
// It returns nil if all handlers accept the visitor without failure.
//
// The order in which handlers are visited is not guaranteed.
func (s HandlerSet) AcceptVisitor(ctx context.Context, v Visitor) error {
	for _, h := range s {
		if err := h.AcceptVisitor(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

// RichHandlerSet is a collection of rich handlers.
type RichHandlerSet map[Identity]RichHandler

// NewRichHandlerSet returns a RichHandlerSet containing the given handlers.
//
// It panics if any of the handler identities conflict.
func NewRichHandlerSet(handlers ...RichHandler) RichHandlerSet {
	s := RichHandlerSet{}

	for _, h := range handlers {
		if !s.Add(h) {
			panic("handler set contains conflicting identities")
		}
	}

	return s
}

// Add adds a handler to the set.
//
// It returns true if the handler was added, or false if the set already
// contained a handler with the same name or key as h.
func (s RichHandlerSet) Add(h RichHandler) bool {
	i := h.Identity()
	for x := range s {
		if i.ConflictsWith(x) {
			return false
		}
	}

	s[i] = h
	return true
}

// Has returns true if s contains h.
func (s RichHandlerSet) Has(h RichHandler) bool {
	x, ok := s[h.Identity()]
	return ok && x == h
}

// ByIdentity returns the handler with the given identity.
func (s RichHandlerSet) ByIdentity(i Identity) (RichHandler, bool) {
	h, ok := s[i]
	return h, ok
}

// ByName returns the handler with the given name.
func (s RichHandlerSet) ByName(n string) (RichHandler, bool) {
	for i, h := range s {
		if i.Name == n {
			return h, true
		}
	}

	return nil, false
}

// ByKey returns the handler with the given key.
func (s RichHandlerSet) ByKey(k string) (RichHandler, bool) {
	for i, h := range s {
		if i.Key == k {
			return h, true
		}
	}

	return nil, false
}

// ByType returns the subset of handlers of the given type.
func (s RichHandlerSet) ByType(t HandlerType) RichHandlerSet {
	subset := RichHandlerSet{}

	for i, h := range s {
		if h.HandlerType() == t {
			subset[i] = h
		}
	}

	return subset
}

// ConsumersOf returns the subset of handlers that consume messages of the given
// type.
func (s RichHandlerSet) ConsumersOf(t message.Type) RichHandlerSet {
	subset := RichHandlerSet{}

	for i, h := range s {
		if h.MessageTypes().Consumed.Has(t) {
			subset[i] = h
		}
	}

	return subset
}

// ProducersOf returns the subset of handlers that produce messages of the given
// type.
func (s RichHandlerSet) ProducersOf(t message.Type) RichHandlerSet {
	subset := RichHandlerSet{}

	for i, h := range s {
		if h.MessageTypes().Produced.Has(t) {
			subset[i] = h
		}
	}

	return subset
}

// IsEqual returns true if o contains the same handlers as s.
func (s RichHandlerSet) IsEqual(o RichHandlerSet) bool {
	if len(s) != len(o) {
		return false
	}

	for i, h := range s {
		x, ok := o[i]
		if !ok || !IsEqual(x, h) {
			return false
		}
	}

	return true
}

// Find returns a handler from the set for which the given predicate function
// returns true.
func (s RichHandlerSet) Find(fn func(RichHandler) bool) (RichHandler, bool) {
	for _, h := range s {
		if fn(h) {
			return h, true
		}
	}

	return nil, false
}

// Filter returns the subset of handlers for which the given predicate function
// returns true.
func (s RichHandlerSet) Filter(fn func(RichHandler) bool) RichHandlerSet {
	subset := RichHandlerSet{}

	for i, h := range s {
		if fn(h) {
			subset[i] = h
		}
	}

	return subset
}

// AcceptRichVisitor visits each handler in the set.
//
// It returns the error returned by the first handler to return a non-nil error.
// It returns nil if all handlers accept the visitor without failure.
//
// The order in which handlers are visited is not guaranteed.
func (s RichHandlerSet) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	for _, h := range s {
		if err := h.AcceptRichVisitor(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

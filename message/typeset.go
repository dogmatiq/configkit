package message

import (
	"github.com/dogmatiq/dogma"
)

// TypeSet is an implementation of TypeCollection that contains a distinct set
// of message types.
type TypeSet map[Type]struct{}

// NewTypeSet returns a TypeSet containing the given types.
func NewTypeSet(types ...Type) TypeSet {
	s := TypeSet{}

	for _, t := range types {
		s[t] = struct{}{}
	}

	return s
}

// TypesOf returns a type set containing the types of the given messages.
func TypesOf(messages ...dogma.Message) TypeSet {
	s := TypeSet{}

	for _, m := range messages {
		s[TypeOf(m)] = struct{}{}
	}

	return s
}

// IntersectionT returns the intersection of the given collections.
//
// That is, it returns a set containing only those types that are present in all
// of the given collections. It returns an empty set if no collections are
// given.
//
// See https://en.wikipedia.org/wiki/Intersection_(set_theory).
func IntersectionT(collections ...TypeCollection) TypeSet {
	r := TypeSet{}

	if len(collections) > 0 {
		collections[0].Each(func(t Type) bool {
			for _, c := range collections[1:] {
				if !c.Has(t) {
					return true
				}
			}

			r.Add(t)

			return true
		})
	}

	return r
}

// UnionT returns the union of the given collections.
//
// That is, it returns a set containing all of the types present in all of the
// given collections. It returns an empty set if no collections are given.
//
// See https://en.wikipedia.org/wiki/Union_(set_theory).
func UnionT(collections ...TypeCollection) TypeSet {
	r := TypeSet{}

	for _, c := range collections {
		c.Each(func(t Type) bool {
			r.Add(t)
			return true
		})
	}

	return r
}

// DiffT returns the difference between a and b.
//
// That is, it returns a set containing all of the types present in a that are
// not present in b.
//
// See https://en.wikipedia.org/wiki/Complement_(set_theory)#Relative_complement.
func DiffT(a, b TypeCollection) TypeSet {
	r := TypeSet{}

	a.Each(func(t Type) bool {
		if !b.Has(t) {
			r.Add(t)
		}

		return true
	})

	return r
}

// Has returns true if s contains t.
func (s TypeSet) Has(t Type) bool {
	_, ok := s[t]
	return ok
}

// HasM returns true if s contains TypeOf(m).
func (s TypeSet) HasM(m dogma.Message) bool {
	return s.Has(TypeOf(m))
}

// Add adds t to s.
//
// It returns true if the type was added, or false if the set already contained
// the type.
func (s TypeSet) Add(t Type) bool {
	if _, ok := s[t]; ok {
		return false
	}

	s[t] = struct{}{}
	return true
}

// AddM adds TypeOf(m) to s.
//
// It returns true if the type was added, or false if the set already contained
// the type.
func (s TypeSet) AddM(m dogma.Message) bool {
	return s.Add(TypeOf(m))
}

// Remove removes t from s.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (s TypeSet) Remove(t Type) bool {
	if _, ok := s[t]; ok {
		delete(s, t)
		return true
	}

	return false
}

// RemoveM removes TypeOf(m) from s.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (s TypeSet) RemoveM(m dogma.Message) bool {
	return s.Remove(TypeOf(m))
}

// IsEqual returns true if o contains the same types as s.
func (s TypeSet) IsEqual(o TypeSet) bool {
	if len(s) != len(o) {
		return false
	}

	for t := range s {
		if _, ok := o[t]; !ok {
			return false
		}
	}

	return true
}

// Each invokes fn once for each type in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// types in the container.
//
// It returns true if fn returned true for all types.
func (s TypeSet) Each(fn func(Type) bool) bool {
	for t := range s {
		if !fn(t) {
			return false
		}
	}

	return true
}

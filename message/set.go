package message

import (
	"iter"
	"maps"

	"github.com/dogmatiq/dogma"
)

// Set is an unorded Set of unique T values.
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet returns a new set containing the given elements.
func NewSet[T comparable](elements ...T) Set[T] {
	var s Set[T]
	s.Add(elements...)
	return s
}

// NamesOf returns a name set containing the names of the given messages.
func NamesOf(messages ...dogma.Message) Set[Name] {
	var names Set[Name]

	for _, m := range messages {
		names.Add(NameOf(m))
	}

	return names
}

// TypesOf returns a type set containing the types of the given messages.
func TypesOf(messages ...dogma.Message) Set[Type] {
	s := Set[Type]{}

	for _, m := range messages {
		s.Add(TypeOf(m))
	}

	return s
}

// Add adds the given elements to the set.
func (s *Set[T]) Add(elements ...T) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}

	for _, e := range elements {
		s.m[e] = struct{}{}
	}
}

// Remove removes the given elements from the set.
func (s *Set[T]) Remove(elements ...T) {
	for _, e := range elements {
		delete(s.m, e)
	}
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.m = nil
}

// Has returns true if the set contains all of the given elements.
func (s Set[T]) Has(elements ...T) bool {
	for _, e := range elements {
		if _, ok := s.m[e]; !ok {
			return false
		}
	}
	return true
}

// All returns an iterator that yields all elements in the set.
func (s Set[T]) All() iter.Seq[T] {
	return maps.Keys(s.m)
}

// Union adds all elements from other to s.
func (s *Set[T]) Union(other Set[T]) {
	for v := range other.m {
		s.Add(v)
	}
}

// IsEqual returns true if s and other contain the same elements.
func (s Set[T]) IsEqual(other Set[T]) bool {
	if len(s.m) != len(other.m) {
		return false
	}

	for k := range s.m {
		if _, ok := other.m[k]; !ok {
			return false
		}
	}

	return true
}

// Len returns the number of elements in the set.
func (s Set[T]) Len() int {
	return len(s.m)
}

// Clone returns a clone of the set.
func (s Set[T]) Clone() Set[T] {
	return Set[T]{m: maps.Clone(s.m)}
}

// Union returns the union of a and b.
func Union[T comparable](a, b Set[T]) Set[T] {
	var u Set[T]
	u.Union(a)
	u.Union(b)
	return u
}

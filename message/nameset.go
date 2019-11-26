package message

import (
	"github.com/dogmatiq/dogma"
)

// NameSet is an implementation of NameCollection that contains a distinct set
// of message names.
type NameSet map[Name]struct{}

// NewNameSet returns a NameSet containing the given names.
func NewNameSet(names ...Name) NameSet {
	s := NameSet{}

	for _, n := range names {
		s[n] = struct{}{}
	}

	return s
}

// NamesOf returns a name set containing the names of the given messages.
func NamesOf(messages ...dogma.Message) NameSet {
	s := NameSet{}

	for _, m := range messages {
		s[NameOf(m)] = struct{}{}
	}

	return s
}

// Has returns true if s contains n.
func (s NameSet) Has(n Name) bool {
	_, ok := s[n]
	return ok
}

// HasM returns true if s contains NameOf(m).
func (s NameSet) HasM(m dogma.Message) bool {
	return s.Has(NameOf(m))
}

// Add adds n to s.
//
// It returns true if the name was added, or false if the set already contained
// the name.
func (s NameSet) Add(n Name) bool {
	if _, ok := s[n]; ok {
		return false
	}

	s[n] = struct{}{}
	return true
}

// AddM adds NameOf(m) to s.
//
// It returns true if the name was added, or false if the set already contained
// the name.
func (s NameSet) AddM(m dogma.Message) bool {
	return s.Add(NameOf(m))
}

// Remove removes n from s.
//
// It returns true if the name was removed, or false if the set did not contain
// the name.
func (s NameSet) Remove(n Name) bool {
	if _, ok := s[n]; ok {
		delete(s, n)
		return true
	}

	return false
}

// RemoveM removes NameOf(m) from s.
//
// It returns true if the name was removed, or false if the set did not contain
// the name.
func (s NameSet) RemoveM(m dogma.Message) bool {
	return s.Remove(NameOf(m))
}

// Each invokes fn once for each name in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// names in the container.
//
// It returns true if fn returned true for all names.
func (s NameSet) Each(fn func(Name) bool) bool {
	for n := range s {
		if !fn(n) {
			return false
		}
	}

	return true
}

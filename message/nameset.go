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

// IntersectionN returns the intersection of the given collections.
//
// That is, it returns a set containing only those names that are present in all
// of the given collections. It returns an empty set if no collections are
// given.
//
// See https://en.wikipedia.org/wiki/Intersection_(set_theory).
func IntersectionN(collections ...NameCollection) NameSet {
	r := NameSet{}

	if len(collections) > 0 {
		collections[0].Range(func(n Name) bool {
			for _, c := range collections[1:] {
				if !c.Has(n) {
					return true
				}
			}

			r.Add(n)

			return true
		})
	}

	return r
}

// UnionN returns the union of the given collections.
//
// That is, it returns a set containing all of the names present in all of the
// given collections. It returns an empty set if no collections are given.
//
// See https://en.wikipedia.org/wiki/Union_(set_theory).
func UnionN(collections ...NameCollection) NameSet {
	r := NameSet{}

	for _, c := range collections {
		c.Range(func(n Name) bool {
			r.Add(n)
			return true
		})
	}

	return r
}

// DiffN returns the difference between a and b.
//
// That is, it returns a set containing all of the names present in a that are
// not present in b.
//
// See https://en.wikipedia.org/wiki/Complement_(set_theory)#Relative_complement.
func DiffN(a, b NameCollection) NameSet {
	r := NameSet{}

	a.Range(func(n Name) bool {
		if !b.Has(n) {
			r.Add(n)
		}

		return true
	})

	return r
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

// IsEqual returns true if o contains the same names as s.
func (s NameSet) IsEqual(o NameSet) bool {
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

// Len returns the number of names in the collection.
func (s NameSet) Len() int {
	return len(s)
}

// Range invokes fn once for each name in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// names in the container.
//
// It returns true if fn returned true for all names.
func (s NameSet) Range(fn func(Name) bool) bool {
	for n := range s {
		if !fn(n) {
			return false
		}
	}

	return true
}

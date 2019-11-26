package message

import (
	"github.com/dogmatiq/dogma"
)

// NameRoles is an implementation of NameCollection that maps message types to
// their roles.
type NameRoles map[Name]Role

// Has returns true if nr contains n.
func (nr NameRoles) Has(n Name) bool {
	_, ok := nr[n]
	return ok
}

// HasM returns true if nr contains NameOf(m).
func (nr NameRoles) HasM(m dogma.Message) bool {
	return nr.Has(NameOf(m))
}

// Add maps n to r.
//
// It returns true if the mapping was added, or false if the map already
// contained the name.
func (nr NameRoles) Add(n Name, r Role) bool {
	if _, ok := nr[n]; ok {
		return false
	}

	nr[n] = r
	return true
}

// AddM adds NameOf(m) to nr.
//
// It returns true if the mapping was added, or false if the map already
// contained the name.
func (nr NameRoles) AddM(m dogma.Message, r Role) bool {
	return nr.Add(NameOf(m), r)
}

// Remove removes n from nr.
//
// It returns true if the name was removed, or false if the set did not contain
// the name.
func (nr NameRoles) Remove(n Name) bool {
	if _, ok := nr[n]; ok {
		delete(nr, n)
		return true
	}

	return false
}

// RemoveM removes NameOf(m) from nr.
//
// It returns true if the name was removed, or false if the set did not contain
// the name.
func (nr NameRoles) RemoveM(m dogma.Message) bool {
	return nr.Remove(NameOf(m))
}

// Each invokes fn once for each name in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// types in the container.
//
// It returns true if fn returned true for all types.
func (nr NameRoles) Each(fn func(Name) bool) bool {
	for n := range nr {
		if !fn(n) {
			return false
		}
	}

	return true
}
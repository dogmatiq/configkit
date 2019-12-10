package message

import (
	"github.com/dogmatiq/dogma"
)

// TypeRoles is an implementation of TypeCollection that maps message types to
// their roles.
type TypeRoles map[Type]Role

// Has returns true if tr contains t.
func (tr TypeRoles) Has(t Type) bool {
	_, ok := tr[t]
	return ok
}

// HasM returns true if tr contains TypeOf(m).
func (tr TypeRoles) HasM(m dogma.Message) bool {
	return tr.Has(TypeOf(m))
}

// Add maps t to r.
//
// It returns true if the mapping was added, or false if the map already
// contained the type.
func (tr TypeRoles) Add(t Type, r Role) bool {
	if _, ok := tr[t]; ok {
		return false
	}

	tr[t] = r
	return true
}

// AddM adds TypeOf(m) to tr.
//
// It returns true if the mapping was added, or false if the map already
// contained the type.
func (tr TypeRoles) AddM(m dogma.Message, r Role) bool {
	return tr.Add(TypeOf(m), r)
}

// Remove removes t from tr.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (tr TypeRoles) Remove(t Type) bool {
	if _, ok := tr[t]; ok {
		delete(tr, t)
		return true
	}

	return false
}

// RemoveM removes TypeOf(m) from tr.
//
// It returns true if the type was removed, or false if the set did not contain
// the type.
func (tr TypeRoles) RemoveM(m dogma.Message) bool {
	return tr.Remove(TypeOf(m))
}

// IsEqual returns true if o contains the same types with the same roles as s.
func (tr TypeRoles) IsEqual(o TypeRoles) bool {
	if len(tr) != len(o) {
		return false
	}

	for t, r := range tr {
		x, ok := o[t]
		if !ok || x != r {
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
func (tr TypeRoles) Each(fn func(Type) bool) bool {
	for t := range tr {
		if !fn(t) {
			return false
		}
	}

	return true
}

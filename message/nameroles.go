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

// IsEqual returns true if o contains the same names with the same roles as s.
func (nr NameRoles) IsEqual(o NameRoles) bool {
	if len(nr) != len(o) {
		return false
	}

	for n, r := range nr {
		x, ok := o[n]
		if !ok || x != r {
			return false
		}
	}

	return true
}

// Len returns the number of names in the collection.
func (nr NameRoles) Len() int {
	return len(nr)
}

// Range invokes fn once for each name in the container.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// names in the container.
//
// It returns true if fn returned true for all names.
func (nr NameRoles) Range(fn func(Name) bool) bool {
	for n := range nr {
		if !fn(n) {
			return false
		}
	}

	return true
}

// RangeByRole invokes fn once for each name in the container that maps to the
// given role.
//
// Iteration stops when fn returns false or once fn has been invoked for all
// names in the container.
//
// It returns true if fn returned true for all names.
func (nr NameRoles) RangeByRole(
	r Role,
	fn func(Name) bool,
) bool {
	for n, xr := range nr {
		if xr == r {
			if !fn(n) {
				return false
			}
		}
	}

	return true
}

// FilterByRole returns the subset of names that have the given role.
func (nr NameRoles) FilterByRole(r Role) NameRoles {
	subset := NameRoles{}

	for n, xr := range nr {
		if r == xr {
			subset[n] = r
		}
	}

	return subset
}

package message

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/collection/sets"
)

// NamesOf returns a name set containing the names of the given messages.
func NamesOf(messages ...dogma.Message) sets.Set[Name] {
	var names sets.Set[Name]

	for _, m := range messages {
		names.Add(NameOf(m))
	}

	return names
}

// TypesOf returns a type set containing the types of the given messages.
func TypesOf(messages ...dogma.Message) sets.Set[Type] {
	var types sets.Set[Type]

	for _, m := range messages {
		types.Add(TypeOf(m))
	}

	return types
}

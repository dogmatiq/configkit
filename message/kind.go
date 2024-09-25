package message

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Kind is an enumeration of the different kinds of messages.
//
// It is similar to a [Role], however it is not tied to a specific application.
type Kind int

const (
	// CommandKind is a [dogma.Message] that also implements [dogma.Command].
	CommandKind Kind = iota

	// EventKind is a [dogma.Message] that also implements [dogma.Event].
	EventKind

	// TimeoutKind is a [dogma.Message] that also implements [dogma.Timeout].
	TimeoutKind
)

func (k Kind) String() string {
	switch k {
	case CommandKind:
		return "command"
	case EventKind:
		return "event"
	case TimeoutKind:
		return "timeout"
	default:
		panic("invalid kind")
	}
}

// KindFor returns the [Kind] of the message with type T.
//
// It panics if T does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func KindFor[T dogma.Message]() Kind {
	var m T
	return KindOf(m)
}

// KindOf returns the [Kind] of m.
func KindOf(m dogma.Message) Kind {
	switch m.(type) {
	case dogma.Command:
		return CommandKind
	case dogma.Event:
		return EventKind
	case dogma.Timeout:
		return TimeoutKind
	default:
		panic(fmt.Sprintf("%T does not implement dogma.Command, dogma.Event or dogma.Timeout", m))
	}
}

func kindFromReflect(rt reflect.Type) Kind {
	if rt.Implements(commandInterface) {
		return CommandKind
	}

	if rt.Implements(eventInterface) {
		return EventKind
	}

	if rt.Implements(timeoutInterface) {
		return TimeoutKind
	}

	panic(fmt.Sprintf("%s does not implement dogma.Command, dogma.Event or dogma.Timeout", rt))
}

var (
	messageInterface = reflect.TypeFor[dogma.Message]()
	commandInterface = reflect.TypeFor[dogma.Command]()
	eventInterface   = reflect.TypeFor[dogma.Event]()
	timeoutInterface = reflect.TypeFor[dogma.Timeout]()
)

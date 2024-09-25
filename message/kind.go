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
	commandInterface = reflect.TypeFor[dogma.Command]()
	eventInterface   = reflect.TypeFor[dogma.Event]()
	timeoutInterface = reflect.TypeFor[dogma.Timeout]()
)

// Switch invokes one of the provided functions based on the [Kind] of m.
//
// It panics if m does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func Switch(
	m dogma.Message,
	command func(dogma.Command),
	event func(dogma.Event),
	timeout func(dogma.Timeout),
) {
	switch m := m.(type) {
	case dogma.Command:
		command(m)
	case dogma.Event:
		event(m)
	case dogma.Timeout:
		timeout(m)
	default:
		panic(fmt.Sprintf(
			"%T implements dogma.Message, but it does not implement dogma.Command, dogma.Event or dogma.Timeout",
			m,
		))
	}
}

// Map invokes one of the provided functions based on the [Kind] of m, and
// returns the result.
//
// It panics if m does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func Map[T any](
	m dogma.Message,
	command func(dogma.Command) T,
	event func(dogma.Event) T,
	timeout func(dogma.Timeout) T,
) (result T) {
	Switch(
		m,
		func(m dogma.Command) { result = command(m) },
		func(m dogma.Event) { result = event(m) },
		func(m dogma.Timeout) { result = timeout(m) },
	)

	return result
}

// TryMap invokes one of the provided functions based on the [Kind] of m, and
// returns the result.
//
// It panics if m does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func TryMap[T any](
	m dogma.Message,
	command func(dogma.Command) (T, error),
	event func(dogma.Event) (T, error),
	timeout func(dogma.Timeout) (T, error),
) (result T, err error) {
	Switch(
		m,
		func(m dogma.Command) { result, err = command(m) },
		func(m dogma.Event) { result, err = event(m) },
		func(m dogma.Timeout) { result, err = timeout(m) },
	)

	return result, err
}

// KindSwitch invokes one of the provided functions based on k.
func KindSwitch(
	k Kind,
	command func(),
	event func(),
	timeout func(),
) {
	switch k {
	case CommandKind:
		command()
	case EventKind:
		event()
	case TimeoutKind:
		timeout()
	default:
		panic("invalid kind")
	}
}

// KindMap invokes one of the provided functions based on k, and returns the
// result.
func KindMap[T any](
	k Kind,
	command func() T,
	event func() T,
	timeout func() T,
) (result T) {
	KindSwitch(
		k,
		func() { result = command() },
		func() { result = event() },
		func() { result = timeout() },
	)

	return result
}

// TryKindMap invokes one of the provided functions based on k, and returns the
// result.
func TryKindMap[T any](
	k Kind,
	command func() (T, error),
	event func() (T, error),
	timeout func() (T, error),
) (result T, err error) {
	KindSwitch(
		k,
		func() { result, err = command() },
		func() { result, err = event() },
		func() { result, err = timeout() },
	)

	return result, err
}

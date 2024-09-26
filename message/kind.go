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

// Symbol returns a character that identifies the message kind when displaying
// message types.
func (k Kind) Symbol() string {
	return MapKind(k, "?", "!", "@")
}

func (k Kind) String() string {
	return MapKind(k, "command", "event", "timeout")
}

// KindFor returns the [Kind] of the message with type T.
//
// It panics if T does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func KindFor[T dogma.Message]() Kind {
	return kindFromReflect(
		reflect.TypeFor[T](),
	)
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
// It provides a compile-time guarantee that all kinds are handled, even if new
// [Kind] values are added in the future.
//
// It panics with a meaningful message if the function associated with m's kind
// is nil.
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
		if command == nil {
			panic("no case function was provided for dogma.Command messages")
		}
		command(m)
	case dogma.Event:
		if event == nil {
			panic("no case function was provided for dogma.Event messages")
		}
		event(m)
	case dogma.Timeout:
		if timeout == nil {
			panic("no case function was provided for dogma.Timeout messages")
		}
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
// It provides a compile-time guarantee that all kinds are handled, even if new
// [Kind] values are added in the future.
//
// It panics with a meaningful message if the function associated with m's kind
// is nil.
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
		mapFunc(command, &result),
		mapFunc(event, &result),
		mapFunc(timeout, &result),
	)

	return result
}

func mapFunc[K dogma.Message, T any](fn func(K) T, v *T) func(K) {
	if fn == nil {
		return nil
	}
	return func(m K) {
		*v = fn(m)
	}
}

// TryMap invokes one of the provided functions based on the [Kind] of m, and
// returns the result and error.
//
// It provides a compile-time guarantee that all kinds are handled, even if new
// [Kind] values are added in the future.
//
// It panics with a meaningful message if the function associated with m's kind
// is nil.
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
		tryMapFunc(command, &result, &err),
		tryMapFunc(event, &result, &err),
		tryMapFunc(timeout, &result, &err),
	)

	return result, err
}

func tryMapFunc[K dogma.Message, T any](
	fn func(K) (T, error),
	v *T,
	err *error,
) func(K) {
	if fn == nil {
		return nil
	}
	return func(m K) {
		*v, *err = fn(m)
	}
}

// SwitchKind invokes one of the provided functions based on k.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [Kind] values are added in the future.
//
// It panics with a meaningful message if the function associated with k.
//
// It panics if k is not a valid [Kind].
func SwitchKind(
	k Kind,
	command func(),
	event func(),
	timeout func(),
) {
	fn := MapKind(k, command, event, timeout)

	if fn == nil {
		panic(fmt.Sprintf("no case function was provided for the %q kind", k))
	}

	fn()
}

// MapKind maps k to a value of type T.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [Kind] values are added in the future.
//
// It panics if k is not a valid [Kind].
func MapKind[T any](k Kind, command, event, timeout T) (result T) {
	switch k {
	case CommandKind:
		return command
	case EventKind:
		return event
	case TimeoutKind:
		return timeout
	default:
		panic("invalid kind")
	}
}

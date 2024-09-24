package message

import (
	"fmt"

	"github.com/dogmatiq/dogma"
)

// Switch invokes one of the provided functions based on the type of m.
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

// Map invokes one of the provided functions based on the type of m and returns
// the result.
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

// TryMap invokes one of the provided functions based on the type of m and
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

package message_test

import (
	"errors"
	"testing"

	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestKindOf(t *testing.T) {
	cases := []struct {
		Message dogma.Message
		Want    Kind
	}{
		{CommandA1, CommandKind},
		{EventA1, EventKind},
		{TimeoutA1, TimeoutKind},
	}

	for _, c := range cases {
		got := KindOf(c.Message)
		if got != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
		}
	}

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		KindOf(nil)
	})
}

func TestKindFor(t *testing.T) {
	if want, got := CommandKind, KindFor[CommandStub[TypeA]](); got != want {
		t.Fatalf("unexpected result: got %q, want %q", got, want)
	}

	if want, got := EventKind, KindFor[EventStub[TypeA]](); got != want {
		t.Fatalf("unexpected result: got %q, want %q", got, want)
	}

	if want, got := TimeoutKind, KindFor[TimeoutStub[TypeA]](); got != want {
		t.Fatalf("unexpected result: got %q, want %q", got, want)
	}

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		KindFor[dogma.Message]()
	})
}

func TestKind(t *testing.T) {
	t.Run("symbol", func(t *testing.T) {
		cases := []struct {
			Kind Kind
			Want string
		}{
			{CommandKind, "?"},
			{EventKind, "!"},
			{TimeoutKind, "@"},
		}

		for _, c := range cases {
			got := c.Kind.Symbol()
			if got != c.Want {
				t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
			}
		}
	})

	t.Run("string", func(t *testing.T) {
		cases := []struct {
			Kind Kind
			Want string
		}{
			{CommandKind, "command"},
			{EventKind, "event"},
			{TimeoutKind, "timeout"},
		}

		for _, c := range cases {
			got := c.Kind.String()
			if got != c.Want {
				t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
			}
		}
	})
}

func TestSwitch(t *testing.T) {
	t.Run("when the message is a command", func(t *testing.T) {
		var message dogma.Command

		Switch(
			CommandA1,
			func(m dogma.Command) { message = m },
			func(m dogma.Event) { t.Fatal("unexpected call to event case") },
			func(m dogma.Timeout) { t.Fatal("unexpected call to timeout case") },
		)

		if message != CommandA1 {
			t.Fatal("command case was not called with the expected message")
		}
	})

	t.Run("when the message is an event", func(t *testing.T) {
		var message dogma.Event

		Switch(
			EventA1,
			func(m dogma.Command) { t.Fatal("unexpected call to command case") },
			func(m dogma.Event) { message = m },
			func(m dogma.Timeout) { t.Fatal("unexpected call to timeout case") },
		)

		if message != EventA1 {
			t.Fatal("event case was not called with the expected message")
		}
	})

	t.Run("when the message is a timeout", func(t *testing.T) {
		var message dogma.Timeout

		Switch(
			TimeoutA1,
			func(m dogma.Command) { t.Fatal("unexpected call to command case") },
			func(m dogma.Event) { t.Fatal("unexpected call to event case") },
			func(m dogma.Timeout) { message = m },
		)

		if message != TimeoutA1 {
			t.Fatal("timeout case was not called with the expected message")
		}
	})

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		Switch(
			nil,
			func(m dogma.Command) { t.Fatal("unexpected call to command case") },
			func(m dogma.Event) { t.Fatal("unexpected call to event case") },
			func(m dogma.Timeout) { t.Fatal("unexpected call to timeout case") },
		)
	})
}

func TestMap(t *testing.T) {
	cases := []struct {
		Message dogma.Message
		Want    string
	}{
		{CommandA1, "command"},
		{EventA1, "event"},
		{TimeoutA1, "timeout"},
	}

	for _, c := range cases {
		got := Map(
			c.Message,
			func(m dogma.Command) string { return "command" },
			func(m dogma.Event) string { return "event" },
			func(m dogma.Timeout) string { return "timeout" },
		)

		if got != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
		}
	}

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		Map(
			nil,
			func(m dogma.Command) string { t.Fatal("unexpected call to command case"); return "" },
			func(m dogma.Event) string { t.Fatal("unexpected call to event case"); return "" },
			func(m dogma.Timeout) string { t.Fatal("unexpected call to timeout case"); return "" },
		)
	})
}

func TestTryMap(t *testing.T) {
	cases := []struct {
		Message dogma.Message
		Want    string
	}{
		{CommandA1, "command"},
		{EventA1, "event"},
		{TimeoutA1, "timeout"},
	}

	for _, c := range cases {
		got, gotErr := TryMap(
			c.Message,
			func(m dogma.Command) (string, error) { return "command", errors.New("command") },
			func(m dogma.Event) (string, error) { return "event", errors.New("event") },
			func(m dogma.Timeout) (string, error) { return "timeout", errors.New("timeout") },
		)

		if got != c.Want {
			t.Fatalf("unexpected result: got %q, want %q", got, c.Want)
		}

		if gotErr == nil {
			t.Fatal("expected an error")
		}

		if gotErr.Error() != c.Want {
			t.Fatalf("unexpected error: got %q, want %q", gotErr, c.Want)
		}
	}

	t.Run("it panics if the message does not implement any of the more specific interfaces", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected a panic")
			}
		}()

		TryMap(
			nil,
			func(m dogma.Command) (string, error) { t.Fatal("unexpected call to command case"); return "", nil },
			func(m dogma.Event) (string, error) { t.Fatal("unexpected call to event case"); return "", nil },
			func(m dogma.Timeout) (string, error) { t.Fatal("unexpected call to timeout case"); return "", nil },
		)
	})
}

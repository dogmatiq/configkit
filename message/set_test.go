package message_test

import (
	"testing"

	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/collection/sets"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestNamesOf(t *testing.T) {
	got := NamesOf(CommandA1, EventA1)
	want := sets.New(NameOf(CommandA1), NameOf(EventA1))

	if !got.IsEqual(want) {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}

func TestTypesOf(t *testing.T) {
	got := TypesOf(CommandA1, EventA1)
	want := sets.New(TypeOf(CommandA1), TypeOf(EventA1))

	if !got.IsEqual(want) {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}

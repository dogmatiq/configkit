package message_test

import (
	"slices"
	"testing"

	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestNamesOf(t *testing.T) {
	got := NamesOf(CommandA1, EventA1)
	want := NewSet(NameOf(CommandA1), NameOf(EventA1))

	if !got.IsEqual(want) {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}

func TestTypesOf(t *testing.T) {
	got := TypesOf(CommandA1, EventA1)
	want := NewSet(TypeOf(CommandA1), TypeOf(EventA1))

	if !got.IsEqual(want) {
		t.Fatalf("unexpected result: got %v, want %v", got, want)
	}
}

func TestSet(t *testing.T) {
	var set Set[int]

	if set.Has(100) {
		t.Fatal("did not expect set to have element")
	}

	set.Add(100)

	if !set.Has(100) {
		t.Fatal("expected set to have element")
	}

	if set.IsEqual(NewSet(200)) {
		t.Fatal("did not expect set to equal a disjoint set with the same number of elements")
	}

	set.Remove(100)

	if set.Has(100) {
		t.Fatal("did not expect set to have element")
	}

	set.Add(100, 200, 300)

	want := 100
	for _, e := range slices.Sorted(set.All()) {
		if e != want {
			t.Fatalf("unexpected element from iterator: got %d, want %d", e, want)
		}
		want += 100
	}

	snapshot := set.Clone()
	set.Union(NewSet(200, 300, 400))

	if set.IsEqual(snapshot) {
		t.Fatalf("union was not performed in-place")
	}

	wantSet := NewSet(100, 200, 300, 400)
	if !set.IsEqual(wantSet) {
		t.Fatalf("union did not produce the correct elements: got %v, want %v", set, wantSet)
	}

	if !set.Has(400) {
		t.Fatal("expected set to have element new element from unioned set")
	}

	if set.Len() != wantSet.Len() {
		t.Fatalf("unexpected number of elements after union: got %d, want %d", set.Len(), wantSet.Len())
	}

	set.Clear()

	if set.Len() != 0 {
		t.Fatalf("expected set to be empty after clearing")
	}
}

func TestUnion(t *testing.T) {
	a := NewSet(100, 200, 300)
	b := NewSet(200, 300, 400)

	union := Union(a, b)

	want := NewSet(100, 200, 300, 400)
	if !union.IsEqual(want) {
		t.Fatalf("unexpected union result: got %v, want %v", a, want)
	}

	if !a.IsEqual(NewSet(100, 200, 300)) {
		t.Fatalf("union modified set a")
	}

	if !b.IsEqual(NewSet(200, 300, 400)) {
		t.Fatalf("union modified set b")
	}
}

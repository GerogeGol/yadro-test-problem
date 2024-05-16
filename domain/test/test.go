package test

import (
	"errors"
	"testing"

	"github.com/GerogeGol/yadro-test-problem/domain/store"
)

var DummyDayTime = store.NewDayTime(0, 0)
var DummyClient = "client"
var DummyTableNumber = 1

func AssertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error at calling arrival: %q", err)
	}
}

func AssertError(t testing.TB, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func AssertTrue(t testing.TB, got bool) {
	t.Helper()
	if !got {
		t.Fatalf("got %t, want 'true'", got)
	}
}

func AssertFalse(t testing.TB, got bool) {
	t.Helper()
	if got {
		t.Fatalf("got %t, want 'false'", got)
	}
}

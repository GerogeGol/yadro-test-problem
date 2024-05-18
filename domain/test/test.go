package test

import (
	"errors"
	"testing"

	"github.com/GerogeGol/yadro-test-problem/domain/service/event"
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
		t.Fatalf("unexpected error: %q", err)
	}
}

func AssertNotNilError(t testing.TB, got error) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected to get any error, but got nil")
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

func AssertEmptyEvent(t testing.TB, e event.Event) {
	t.Helper()
	if !event.IsEmpty(e) {
		t.Fatalf("got error event: %#v", e)
	}
}

func AssertNoErrorEvent(t testing.TB, e event.Event) {
	t.Helper()
	err, ok := e.(*event.ErrorEvent)
	if ok {
		t.Fatalf("got error event: %q", err.Err())
	}
}

func AssertErrorEvent(t testing.TB, e event.Event, err error) {
	t.Helper()
	errEvent, ok := e.(*event.ErrorEvent)
	if !ok {
		t.Fatalf("expected error event got: %#v", e)
	}
	AssertError(t, errEvent.Err(), err)
}

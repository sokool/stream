package stream_test

import (
	"testing"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

func TestNewEvent(t *testing.T) {
	id, err := stream.ParseRootID("ja5285.Thread")
	if err != nil {
		t.Fatal(err)
	}
	m := threads.ThreadStarted{Moderator: "Tom", Channel: "#general"}
	e, err := stream.NewEvent(id, m, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !e.Belongs(id) {
		t.Fatal("ble")
	}
	if e.IsEmpty() {
		t.Fatal("expected nonempty event")
	}
	if e.Root() != id.Type() {
		t.Fatal()
	}
	if e.Type() != "Started" {
		t.Fatal()
	}
	if e.Name() != "ThreadStarted" {
		t.Fatal()
	}
	if e.CreatedAt().IsZero() {
		t.Fatal()
	}
	if e.Sequence() != 1 {
		t.Fatal()
	}
	if e.Stream() != id.ID() {
		t.Fatal()
	}
}

func TestNewEvents(t *testing.T) {
	r, err := stream.NewAggregate("Uh3D9L13", threads.New, threads.Events)
	if err != nil {
		t.Fatal(err)
	}

	s := func(t *threads.Thread) error { return t.Start("#general", "tom") }

	if err = r.Run(s); err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	if s := r.Events().Size(); s != 2 {
		t.Fatalf("expected 2 events, got %d", s)
	}

	if r.Events()[0].String() != "Uh3D9L13:1:Thread[Started]" {
		t.Fatal()
	}

	if r.Events()[1].String() != "Uh3D9L13:2:Thread[Joined]" {
		t.Fatal()
	}
}

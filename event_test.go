package stream_test

import (
	"testing"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

//func TestEvent_String(t *testing.T) {
//	type scenario struct {
//		description string
//		event       stream.Event
//		err         bool
//	}
//
//	cases := []scenario{
//		{"my first subtest", stream.NewEvent()},
//	}
//
//	for _, c := range cases {
//		t.Run(c.description, func(t *testing.T) {
//			_, err :=
//			if c.err && err == nil {
//				t.Fatalf("error expected")
//			}
//			if !c.err && err != nil {
//				t.Fatalf("no error expected, got %v", err)
//			}
//		})
//	}
//}

func TestNewEvent(t *testing.T) {
	s, err := stream.ParseSequence("ja5285.Thread")
	if err != nil {
		t.Fatal(err)
	}
	id := s.ID()

	e, err := stream.NewEvent(s, threads.ThreadStarted{})
	if err != nil {
		t.Fatal(err)
	}
	if !e.Belongs(id) {
		t.Fatal("ble")
	}
	if e.IsEmpty() {
		t.Fatal("expected nonempty event")
	}
	if e.Stream() != id {
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
	if e.Sequence() != 0 {
		t.Fatal()
	}
	if e.Stream() != id {
		t.Fatal()
	}
	if n := e.String(); n != "5776c729:Thread[Started]" {
		t.Fatalf("expected 5776c729:Thread[Started] got %s", n)
	}

}

func TestNewEvents(t *testing.T) {
	e, err := stream.NewEvents(
		stream.MustSequence[threads.Thread]("Uh3D9L13"),
		threads.ThreadStarted{},
		threads.ThreadJoined{},
		threads.ThreadClosed{},
	)

	if err != nil {
		t.Fatal(err)
	}

	if n := e.Size(); n != 3 {
		t.Fatalf("expected 3 events, got %d", n)
	}

	if s := e[0].String(); s != "497b882d:Thread[Started]#1" {
		t.Fatalf("expected 497b882d:Thread[Started]#1 got %s", s)
	}

	if s := e[1].String(); s != "497b882d:Thread[Joined]#2" {
		t.Fatalf("expected 497b882d:Thread[Started]#1 got %s", s)
	}

	if s := e[2].String(); s != "497b882d:Thread[Closed]#3" {
		t.Fatalf("expected 497b882d:Thread[Started]#1 got %s", s)
	}
}

// Uh3D9L1fa:Thread
// Uh3D9L1fa:Thread#3
// Uh3D9L1fa:Thread[Started]#3
// Uh3D9L1fa:Thread[Started Joined]#2

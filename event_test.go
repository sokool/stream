package stream_test

import (
	"fmt"
	"testing"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
)

func Test(t *testing.T) {
	c, _ := model.NewThread("48")
	n, _ := stream.NewRootID(c)

	e, err := stream.NewEvent(n, model.ThreadStarted{}, 1)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(e.ID())
	//fmt.Println(e.Sequence())
	//fmt.Println(e.Type())
	//fmt.Println(e.Body())
	//fmt.Println(e.Namespace())
	fmt.Println(e.GoString())

}

func TestNewEvent(t *testing.T) {
	id, err := stream.ParseRootID("ja5285.Thread")
	if err != nil {
		t.Fatal(err)
	}
	m := model.ThreadStarted{Moderator: "Tom", Channel: "#general"}
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

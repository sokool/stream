package stream_test

import (
	"fmt"
	"testing"
)

func TestAggregates(t *testing.T) {
	id := "73HdaUj"
	threads := &Threads{OnCreate: NewThread}

	if err := threads.Execute(id, func(t *Thread) error { return t.Start("#fire-up", "tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	if err := threads.Execute(id, func(t *Thread) error { return t.Message("tom@on.de", "hi there") }); err != nil {
		t.Fatal(err)
	}

	x, err := threads.Read(id)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(threads)
	fmt.Println(x)
}

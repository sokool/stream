package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/example/chat/repository"
	"testing"
)

func TestAggregates(t *testing.T) {
	id, threads := "73HdaUj", repository.NewThreads()
	if err := threads.Execute(id, func(t *model.Thread) error { return t.Start("#fire-up", "tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	if err := threads.Execute(id, func(t *model.Thread) error { return t.Message("tom@on.de", "hi there") }); err != nil {
		t.Fatal(err)
	}

	fmt.Println(threads)

	//x, err := threads.Read(id)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//
	//if x.String() != "73HdaUj.Thread#3" {
	//	t.Fatal()
	//}
}

package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/example/chat/repository"
	"testing"
)

func TestAggregates(t *testing.T) {
	r := repository.NewThreads()
	d := "abc"
	err := r.Execute(d, func(m *model.Thread) error {
		return m.Start("#fire-up", "tom@on.de")
	})
	if err != nil {
		t.Fatal(err)
	}

	err = r.Execute(d, func(m *model.Thread) error {
		return m.Start("#fire-up", "tom@on.de")
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r)
	//n := stream.MustNamespace("lajf89371.Chat")
	//s := stream.NewSequence(n)
	//e := model.ThreadStarted{
	//	Moderator: "tom@on.de",
	//	Channel:   "#fire-up",
	//}
	//m := stream.NewMessage(s, e)
	//fmt.Println(m.GoString())
	//fmt.Println(m.String())

}

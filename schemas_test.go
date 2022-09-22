package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/repository"
	"testing"
)

func TestNewScheme(t *testing.T) {
	threads := repository.NewChats()

	for i := range threads.Threads.Events {
		fmt.Println(threads.Threads.Events[i].String())
	}

	//x := threads.(stream.Aggregate[stream.Root])
	//var y stream.Aggregate[stream.Root]

	//fmt.Println(x)
	//s := stream.NewDomain[stream.Aggregate[stream.Root]]()
	//s.Register(threads.(*stream.Aggregate[stream.Root]))

	//fmt.Println(s.Register(threads))
	//o := []stream.Aggregate[stream.Root]{
	//	stream.Aggregate[*model.Thread]{},
	//}
	//fmt.Println(o)
	//s.Register(x.(stream.Aggregate[stream.Root]))
}

package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/repository"
	"testing"
)

func TestNewScheme(t *testing.T) {
	threads := repository.NewChats()

	for i := range threads.Threads.Events {
		fmt.Println(threads.Threads.Events[i])
	}

	//x := threads.(stream.Aggregates[stream.Root])
	//var y stream.Aggregates[stream.Root]

	//fmt.Println(x)
	//s := stream.NewDomain[stream.Aggregates[stream.Root]]()
	//s.Register(threads.(*stream.Aggregates[stream.Root]))

	//fmt.Println(s.Register(threads))
	//o := []stream.Aggregates[stream.Root]{
	//	stream.Aggregates[*model.Thread]{},
	//}
	//fmt.Println(o)
	//s.Register(x.(stream.Aggregates[stream.Root]))
}

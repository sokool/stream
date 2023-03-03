package stream_test

import (
	"fmt"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat"
)

func ExampleAggregate_Execute() {
	s := stream.New(&stream.Configuration{})
	chats, _ := chat.New(s)
	id := "k8Duq81o"
	chat, err := chats.Threads.Get(id)
	if err != nil {
		return
	}
	fmt.Println(chat)

	if err = chat.Start("elo", "dood"); err != nil {
		return
	}
	fmt.Println(chat)

	if err = chats.Threads.Set(chat); err != nil {
		return
	}

	fmt.Println(chat)

	//x, _ := threads.Members.Store.Load(nil)

	// Output:
	// k8Duq81o.Thread#0
	// k8Duq81o.Thread#0->2
	// k8Duq81o.Thread#2
}

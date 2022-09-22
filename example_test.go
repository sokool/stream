package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/repository"
)

func ExampleAggregate_Execute() {
	id, threads := "k8Duq81o", repository.NewChats()
	chat, err := threads.Threads.Get(id)
	if err != nil {
		return
	}
	fmt.Println(chat)

	if err = chat.Start("elo", "dood"); err != nil {
		return
	}
	fmt.Println(chat)

	if err = threads.Threads.Set(chat); err != nil {
		return
	}

	fmt.Println(chat)

	//x, _ := threads.Members.Store.Load(nil)

	// Output:
	// k8Duq81o.Thread#0
	// k8Duq81o.Thread#0->2
	// k8Duq81o.Thread#2
}

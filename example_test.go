package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/repository"
)

func ExampleAggregate_Execute() {
	id, threads := "k8Duq81o", repository.NewThreads()
	chat, err := threads.Read(id)
	if err != nil {
		return
	}
	fmt.Println(chat)

	if err = chat.Start("elo", "dood"); err != nil {
		return
	}
	fmt.Println(chat)

	if err = threads.Write(chat); err != nil {
		return
	}
	fmt.Println(chat)

	// Output:
	// k8Duq81o.Thread#0
	// k8Duq81o.Thread#0->2
	// k8Duq81o.Thread#2
}

package stream_test

import (
	"context"
	"fmt"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

func ExampleAggregate_Run() {
	id, _ := stream.NewID[threads.Thread]("k8Duq81o")
	chat, err := stream.NewAggregate(id, threads.New, threads.Events)
	if err != nil {
		return
	}
	fmt.Println(chat)

	if err = chat.Run(func(t *threads.Thread) error { return t.Start("elo", "dood") }); err != nil {
		return
	}
	fmt.Println(chat)

	if _, err = chat.WriteTo(stream.MemoryEventStore); err != nil {
		return
	}
	fmt.Println(chat)

	// Output:
	// 049fab7a.Thread
	// 049fab7a.Thread->2
	// 049fab7a.Thread#2
}

type Person struct {
}

func (p *Person) Context() stream.Context {
	return context.TODO()
}

func (p *Person) IsGranted(resource string) error {
	fmt.Println(resource)
	return nil
}

package repository

import (
	"fmt"
	"github.com/sokool/stream"
	. "github.com/sokool/stream/example/chat/model"
)

type Threads = stream.Aggregate[*Thread]

func NewThreads() *Threads {
	return &Threads{
		OnCreate: func(id string) (*Thread, error) {
			t, err := NewThread(id)
			fmt.Println("on.create", t)
			return t, err
		},
		OnRead: func(t *Thread) error {
			fmt.Println("on.read", t)
			return nil
		},
		OnWrite: func(t *Thread) error {
			fmt.Println("on.write", t)
			return nil
		},
		OnCommit: func(t *Thread, e []stream.Event[any]) error {
			fmt.Println("on.commit", t, len(e))
			return nil
		},
		OnCacheCleanup: nil,
	}
}

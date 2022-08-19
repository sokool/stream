package repository

import (
	"github.com/sokool/stream"
	. "github.com/sokool/stream/example/chat/model"
)

type Threads = stream.Aggregate[*Thread, Event]

func NewThreads() *Threads {
	return &Threads{
		OnCreate: func(id string) (*Thread, error) {
			return NewThread(id)
		},
		Events:         nil,
		OnRead:         nil,
		OnWrite:        nil,
		OnChange:       nil,
		OnCacheCleanup: nil,
	}
}

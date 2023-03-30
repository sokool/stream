package repository

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

type Threads = stream.Aggregates[*threads.Thread]

func NewThreads() *Threads {
	return stream.NewAggregates(threads.New, threads.Events)
}

package repository

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

type Threads = stream.Aggregates[*threads.Thread]

func NewThreads(s *stream.Engine) (*Threads, error) {
	return stream.NewAggregates(threads.New, threads.Events).Compose(s), nil

	//a := &Threads{
	//	OnCreate: threads.New,
	//	OnLoad: func(t *threads.Thread) error {
	//		return nil
	//	},
	//	OnCommit: func(t *threads.Thread, e stream.Events) error {
	//		return nil
	//	},
	//	OnSave: func(t *threads.Thread) error {
	//		return nil
	//	},
	//	Events: stream.Schemas{
	//		threads.ThreadStarted{}: {
	//			Description: "thread starts automatically, when there is a longer break between messages",
	//			//Transaction: m,
	//		},
	//		threads.ThreadMessage{}: {
	//			Transaction: "Conversations",
	//		},
	//		threads.ThreadJoined{}: {
	//			Transaction: "Members",
	//		},
	//		threads.ThreadLeft{}: {
	//			Transaction: "Members",
	//		},
	//		threads.ThreadMuted{}: {},
	//		threads.ThreadKicked{}: {
	//			Transaction: "Members",
	//		},
	//	},
	//	OnCacheCleanup:     nil,
	//	CleanCacheAfter:    -1,
	//	LoadEventsInChunks: 8,
	//}

	//return a, a.Compose(s)
}

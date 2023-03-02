package repository

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

type Threads struct {
	*stream.Aggregates[*threads.Thread]
}

func NewThreads() *Threads {
	return &Threads{
		&stream.Aggregates[*threads.Thread]{
			Description: "",
			OnCreate: func(id string) (*threads.Thread, error) {
				t, err := threads.NewThread(id)
				//fmt.Println("on.create", t)
				return t, err
			},
			OnLoad: func(t *threads.Thread) error {
				//fmt.Println("on.load", t)
				return nil
			},
			OnCommit: func(t *threads.Thread, e stream.Events) error {
				//fmt.Println("on.commit", t, len(e))
				return nil
			},
			OnSave: func(t *threads.Thread) error {
				//fmt.Println("on.save", t)
				return nil
			},
			Events: stream.Schemas{
				threads.ThreadStarted{}: {
					Description: "thread starts automatically, when there is a longer break between messages",
					//Transaction: m,
				},
				threads.ThreadMessage{}: {
					Transaction: "Conversations",
				},
				threads.ThreadJoined{}: {
					Transaction: "Members",
				},
				threads.ThreadLeft{}: {
					Transaction: "Members",
				},
				threads.ThreadMuted{}: {},
				threads.ThreadKicked{}: {
					Transaction: "Members",
				},
			},
			OnCacheCleanup:     nil,
			CleanCacheAfter:    -1,
			LoadEventsInChunks: 8,
		},
	}
}

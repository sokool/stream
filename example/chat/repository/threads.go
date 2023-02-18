package repository

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
)

type Threads struct {
	*stream.Aggregates[*model.Thread]
}

func NewThreads() *Threads {
	return &Threads{
		&stream.Aggregates[*model.Thread]{
			Description: "",
			OnCreate: func(id string) (*model.Thread, error) {
				t, err := model.NewThread(id)
				//fmt.Println("on.create", t)
				return t, err
			},
			OnLoad: func(t *model.Thread) error {
				//fmt.Println("on.load", t)
				return nil
			},
			OnCommit: func(t *model.Thread, e stream.Events) error {
				//fmt.Println("on.commit", t, len(e))
				return nil
			},
			OnSave: func(t *model.Thread) error {
				//fmt.Println("on.save", t)
				return nil
			},
			Events: stream.Schemas{
				model.ThreadStarted{}: {
					Description: "thread starts automatically, when there is a longer break between messages",
					//Transaction: m,
				},
				model.ThreadMessage{}: {
					Transaction: "Conversations",
				},
				model.ThreadJoined{}: {
					Transaction: "Members",
				},
				model.ThreadLeft{}: {
					Transaction: "Members",
				},
				model.ThreadMuted{}: {},
				model.ThreadKicked{}: {
					Transaction: "Members",
				},
			},
			OnCacheCleanup:     nil,
			CleanCacheAfter:    -1,
			LoadEventsInChunks: 8,
		},
	}
}

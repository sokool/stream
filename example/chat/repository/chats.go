package repository

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
)

type Chats struct {
	Threads  *stream.Aggregate[*model.Thread]
	Members  *stream.Projection[*Member]
	Messages *stream.Projection[*Messages]
}

func NewChats() *Chats {
	return &Chats{
		Threads: &stream.Aggregate[*model.Thread]{
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
					Coupling:    []string{"Messages", "Member"},
				},
				model.ThreadMessage{}: {},
				model.ThreadJoined{}:  {},
				model.ThreadLeft{}:    {},
				model.ThreadKicked{}:  {},
				model.ThreadMuted{}:   {},
				model.ThreadKicked{}:  {},
			},
			OnCacheCleanup:     nil,
			CleanCacheAfter:    -1,
			LoadEventsInChunks: 8,
		},
		Members: &stream.Projection[*Member]{
			Documents: NewMembers(),
		},
		Messages: &stream.Projection[*Messages]{
			Documents: NewMessagez(),
		},
	}
}

func (t *Chats) Thread(id string, command func(*model.Thread) error) error {
	return t.Threads.Execute(id, command)
}

func (t *Chats) Register(s *stream.Domain) error {
	return s.Register(t.Threads, t.Messages, t.Members)
}

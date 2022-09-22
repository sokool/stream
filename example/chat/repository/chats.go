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
			OnRead: func(t *model.Thread) error {
				//fmt.Println("on.read", t)
				return nil
			},
			OnWrite: func(t *model.Thread) error {
				//fmt.Println("on.write", t)
				return nil
			},
			OnCommit: func(t *model.Thread, e stream.Events) error {
				//fmt.Println("on.commit", t, len(e))
				return nil
			},
			Events: stream.Schemas{
				stream.NewScheme(model.ThreadStarted{}).Name("Started").Description("some desc").Couple(""),
				stream.NewScheme(model.ThreadMessage{}),
				stream.NewScheme(model.ThreadJoined{}),
				stream.NewScheme(model.ThreadLeft{}),
				stream.NewScheme(model.ThreadKicked{}),
				stream.NewScheme(model.ThreadMuted{}),
				stream.NewScheme(model.ThreadClosed{}),
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

func (t *Chats) Membersx() {}

func (t *Chats) Register(s *stream.Domain) error {
	return s.Register(t.Threads, t.Messages, t.Members)
}

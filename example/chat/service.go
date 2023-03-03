package chat

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/repository"
)

type Service struct {
	Stream        *stream.Engine
	Threads       *repository.Threads
	Members       *repository.Members
	Conversations *repository.Conversations
}

func New(se *stream.Engine) (*Service, error) {
	var s = Service{Stream: se}
	var err error
	if s.Threads, err = repository.NewThreads(se); err != nil {
		return nil, err
	}
	if s.Members, err = repository.NewMembers(se); err != nil {
		return nil, err
	}
	if s.Conversations, err = repository.NewConversations(se); err != nil {
		return nil, err
	}
	return &s, nil
}

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
	s := Service{
		Stream:        se,
		Threads:       repository.NewThreads(),
		Members:       repository.NewMembers(),
		Conversations: repository.NewConversations(),
	}
	return &s, se.Compose(s.Threads, s.Members, s.Conversations)
}

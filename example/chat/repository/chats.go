package repository

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/store/mysql"
	"math/rand"
	"os"
	"time"
)

type Chats struct {
	*Threads
	*Members
	*Conversations
}

func NewChats() *Chats {
	return &Chats{
		Threads:       NewThreads(),
		Members:       NewMembers(),
		Conversations: NewConversations(),
	}
}

func (t *Chats) Thread(id string, command func(*model.Thread) error) error {
	return t.Threads.Execute(id, command)
}

func (t *Chats) Register(d *stream.Domain) error {
	return d.Register(t.Threads, t.Conversations, t.Members)
}

func storage[E stream.Entity](fn stream.EntityFunc[E]) (stream.Entities[E], error) {
	if cdn := os.Getenv("MYSQL_EVENT_STORE"); cdn != "" {
		c, err := mysql.NewConnection(cdn, nil)
		if err != nil {
			return nil, err
		}

		m, err := mysql.NewTable[E](c, fn)
		if err != nil {
			return nil, err
		}

		return m, nil
	}

	return stream.NewEntities[E](fn), nil
}

func delay(x time.Duration) {
	max := int64(x)
	min := int64(time.Millisecond * 100)
	t := time.Duration(rand.Int63n(max-min) + min)

	time.Sleep(t)
}

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
	Threads       *Threads
	Members       *Members
	Conversations *Conversations
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

	c := []stream.Component{
		t.Threads,
		&stream.Projection[*Member]{
			Documents: t.Members,
		},
		&stream.Projection[*Conversation]{
			Documents: t.Conversations,
		},
	}

	return d.Register(c...)
}

func storage[E stream.Entity](fn stream.EntityFunc[E]) (stream.Documents[E], error) {
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

	return stream.NewDocuments[E](fn), nil
}

func delay(x time.Duration) {
	max := int64(x)
	min := int64(time.Millisecond * 100)
	t := time.Duration(rand.Int63n(max-min) + min)

	time.Sleep(t)
}

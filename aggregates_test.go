package stream_test

import (
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/example/chat/repository"
	"github.com/sokool/stream/store/mysql"
	"os"
	"testing"
	"time"
)

func TestAggregates(t *testing.T) {
	id, chats := "73HdaUj", repository.NewChats()

	if err := NewDomain(t).Register(chats); err != nil {
		t.Fatal(err)
	}

	type Thread = model.Thread
	if err := chats.Thread(id, func(t *Thread) error { return t.Start("#fire-up", "tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Message("tom@on.de", "hi there") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Join("greg@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
}

func NewDomain(t *testing.T) *stream.Domain {
	return stream.NewDomain(&stream.Configuration{
		EventStore: func(s *stream.Schemas, l stream.Printer) stream.EventStore {
			es, err := mysql.NewEventsStore(os.Getenv("MYSQL_EVENT_STORE"), s, l)
			if err != nil {
				t.Fatal(err)
			}
			return es
		},
	})
}

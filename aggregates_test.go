package stream_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/repository"
	"github.com/sokool/stream/example/chat/threads"
	"github.com/sokool/stream/store/mysql"
)

func TestAggregates(t *testing.T) {
	d := NewEngine(t)
	id, ch, chats := fake.CharactersN(6), "#"+strings.ReplaceAll(strings.ToLower(fake.Street()), " ", "-"), repository.NewChats()

	if err := d.Compose(chats); err != nil {
		t.Fatal(err)
	}

	type Thread = threads.Thread
	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Start(ch, "tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Message("tom@on.de", "hi there") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Join("greg@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Message("greg@gog.pl", "crusher!") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Join("mark@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Message("tom@on.de", "fine, thx!") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Threads.Execute(id, func(t *Thread) error { return t.Leave("tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 7)
}

func NewEngine(t *testing.T) *stream.Engine {
	return stream.New(&stream.Configuration{
		Name: "MyCoolTestDomain",
		EventStore: func(l stream.Printer) stream.EventStore {
			host := os.Getenv("MYSQL_EVENT_STORE")
			if host == "" {
				return stream.NewEventStore()
			}

			es, err := mysql.NewEventsStore(host, l)
			if err != nil {
				t.Fatal(err)
			}
			return es
		},
	})
}

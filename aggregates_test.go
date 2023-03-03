package stream_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat"
	"github.com/sokool/stream/example/chat/threads"
	"github.com/sokool/stream/store/mysql"
)

func TestAggregates(t *testing.T) {
	chats, err := chat.New(NewEngine(t))
	if err != nil {
		t.Fatal(err)
	}

	id, ch := fake.CharactersN(6), "#"+strings.ReplaceAll(strings.ToLower(fake.Street()), " ", "-")
	type Thread = threads.Thread
	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Start(ch, "tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Message("tom@on.de", "hi there") }); err != nil {
		t.Fatal(err)
	}

	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Join("greg@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Message("greg@gog.pl", "crusher!") }); err != nil {
		t.Fatal(err)
	}

	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Join("mark@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Message("tom@on.de", "fine, thx!") }); err != nil {
		t.Fatal(err)
	}

	if err = chats.Threads.Execute(id, func(t *Thread) error { return t.Leave("tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 7)
}

func TestAggregates_SetGet(t *testing.T) {
	id, se := fake.CharactersN(8), NewEngine(t)
	chats, err := chat.New(se)
	t1, err := chats.Threads.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	t2, err := chats.Threads.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if t1 != t2 {
		t.Fatal()
	}
	if t1.ID() != id {
		t.Fatal()
	}
	if t1.Version() != 0 {
		t.Fatal()
	}
	if len(t1.Uncommitted(false)) != 0 {
		t.Fatal()
	}
	if err = t1.Start("#general", "tom"); err != nil {
		t.Fatal(err)
	}
	if t1.Version() != 0 {
		t.Fatal()
	}
	if len(t1.Uncommitted(false)) != 2 {
		t.Fatal()
	}
	if err = chats.Threads.Set(t1); err != nil {
		t.Fatal(err)
	}
	if t1.Version() != 2 {
		t.Fatal()
	}
	if len(t1.Uncommitted(false)) != 0 {
		t.Fatal()
	}
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

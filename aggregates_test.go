package stream_test

import (
	"fmt"
	"github.com/icrowley/fake"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/example/chat/repository"
	"github.com/sokool/stream/store/mysql"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAggregates(t *testing.T) {
	d := NewDomain(t)
	id, ch, chats := fake.CharactersN(6), "#"+strings.ReplaceAll(strings.ToLower(fake.Street()), " ", "-"), repository.NewChats()

	if err := chats.Register(d); err != nil {
		t.Fatal(err)
	}

	type Thread = model.Thread
	if err := chats.Thread(id, func(t *Thread) error { return t.Start(ch, "tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Message("tom@on.de", "hi there") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Join("greg@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Message("greg@gog.pl", "crusher!") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Join("mark@gog.pl") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Message("tom@on.de", "fine, thx!") }); err != nil {
		t.Fatal(err)
	}

	if err := chats.Thread(id, func(t *Thread) error { return t.Leave("tom@on.de") }); err != nil {
		t.Fatal(err)
	}

	fmt.Println("done")
	time.Sleep(time.Second * 7)
}

func NewDomain(t *testing.T) *stream.Domain {
	return stream.NewDomain(&stream.Configuration{
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

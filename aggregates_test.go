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
	id, chats, es := "73HdaUj", repository.NewChats(), NewEventStore(t)
	if err := stream.NewDomain(es).Register(chats); err != nil {
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

func NewEventStore(t *testing.T) stream.EventStore {
	es, err := mysql.NewEventsStore(os.Getenv("MYSQL_EVENT_STORE"), &stream.Schemas{})
	if err != nil {
		t.Fatal(err)
	}
	return es
}

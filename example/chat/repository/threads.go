package repository

import (
	"fmt"
	"github.com/sokool/stream"
	. "github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/store/mysql"
	"os"
)

type Threads = stream.Aggregate[*Thread]

func NewThreads() *Threads {
	es, _ := mysql.NewEventsStore(os.Getenv("MYSQL_EVENT_STORE"), &stream.Schemas{})

	return &Threads{
		Description: "",
		OnCreate: func(id string) (*Thread, error) {
			t, err := NewThread(id)
			fmt.Println("on.create", t)
			return t, err
		},
		OnRead: func(t *Thread) error {
			fmt.Println("on.read", t)
			return nil
		},
		OnWrite: func(t *Thread) error {
			fmt.Println("on.write", t)
			return nil
		},
		OnCommit: func(t *Thread, e stream.Events) error {
			fmt.Println("on.commit", t, len(e))
			return nil
		},
		Events: stream.Schemas{
			stream.NewScheme(ThreadStarted{}).
				Name("Started").
				Description("some desc").
				Couple(""),
			stream.NewScheme(ThreadMessage{}),
			stream.NewScheme(ThreadJoined{}),
			stream.NewScheme(ThreadLeft{}),
			stream.NewScheme(ThreadKicked{}),
			stream.NewScheme(ThreadMuted{}),
			stream.NewScheme(ThreadClosed{}),
		},
		OnCacheCleanup:     nil,
		CleanCacheAfter:    -1,
		LoadEventsInChunks: 8,
		Store:              es,
	}
}

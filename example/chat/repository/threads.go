package repository

import (
	"github.com/sokool/stream"
	. "github.com/sokool/stream/example/chat/model"
)

type Threads = stream.Aggregate[*Thread, Event]

func NewThreads() *Threads {
	return &Threads{
		OnCreate: func(id string) (*Thread, error) {
			return NewThread(id), nil
		},
		Events:         nil,
		OnRead:         nil,
		OnWrite:        nil,
		OnChange:       nil,
		OnCacheCleanup: nil,
	}

}

func a() {

}

//type View struct {
//	Started bool
//	Members model.Participants
//}
//
//func NewAggregate(options ...Option) stream.Aggregate {
//	return stream.Aggregate{
//		Name:        "Chat",
//		Description: "example thread in chat",
//		Events: []stream.Event{
//			stream.NewEvent(ThreadStarted{}),
//			stream.NewEvent(ThreadJoined{}).Description("participant has been joined to make an conversation"),
//			stream.NewEvent(ThreadMessage{}).Description("participant send message"),
//			stream.NewEvent(ThreadLeft{}),
//			stream.NewEvent(ThreadKicked{}).Description("participant has been kicked by moderator"),
//			stream.NewEvent(ThreadMuted{}).Description("moderator "),
//			stream.NewEvent(ThreadClosed{}),
//			stream.NewEvent(recalled{}).Description("internal event, invoked by goes"),
//		},
//		OnCreate: func(n stream.ID) stream.Root {
//			t := New(n.Value())
//			for i := range options {
//				options[i](t)
//			}
//			return t
//		},
//		OnSession: func(r stream.Root) (stream.Session, error) {
//			return nil, nil
//		},
//		OnCommand: func(s stream.Session, r stream.Root) error {
//			//fmt.Println("on.command", sa.(*Thread).changelog.String())
//			return nil
//		},
//		OnChange: func(n stream.ID, m []stream.Message) error {
//			//v := sa.(*Thread).Serialize()
//			//b, _ := json.MarshalIndent(v, "", "\t")
//			//fmt.Println("on.change", sa.(*Thread).changelog.String())
//			return nil
//		},
//		OnCacheCleanup: func(n stream.ID) error {
//			//fmt.Println("cache:", sa.(*Thread).changelog)
//			return nil
//		},
//		CleanCacheAfter: time.Hour,
//		//Log: log.Printf,
//	}
//}
//
//type Option func(m *model.Thread)
//
//func Recall(after time.Duration, random bool, max int) Option {
//	return func(c *model.Thread) {
//		c.recalls.max = max
//		c.recalls.after = after
//		if random {
//			c.recalls.after = time.Duration(rand.Int63n(int64(after)))
//
//		}
//	}
//}
//
//func Delay(d time.Duration) Option {
//	return func(c *model.Thread) {
//		c.settings.delay = time.Duration(rand.Int63n(int64(d)))
//	}
//}
//
//func Expire(after time.Duration, random bool) Option {
//	return func(c *model.Thread) {
//		//if random {
//		//	c.changelog.Evict(time.Duration(rand.Int63n(int64(after))))
//		//	return
//		//}
//		//
//		//c.changelog.Evict(after)
//	}
//}
//
//type Command func(c *macho.Thread) error

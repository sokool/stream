package stream

import (
	"fmt"
	"testing"
	"time"
)

func TestAggregates(t *testing.T) {
	c := Aggregate[*chat]{
		Type: "chat",
		OnCreate: func(id ID) (*chat, error) {
			return newChat(id.String())
		},
	}

	n := NewAggregates(c)
	d := MustID("911.chat")
	err := n.Execute(d, func(c *chat) error {
		fmt.Println("exec", c)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

type chat struct {
	id string
}

func newChat(id string) (*chat, error) {
	return &chat{id}, nil

}

func (r *chat) ID() string {
	return r.id
}

func (r *chat) Name() string {
	return "chat"
}

func (r *chat) Uncommitted(b bool) []Event {
	//TODO implement me
	panic("implement me")
}

func (r *chat) Commit(event Event, time time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (r *chat) String() string {
	return r.id
}

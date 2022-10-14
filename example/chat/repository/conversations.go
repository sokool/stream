package repository

import (
	"fmt"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"strings"
	"time"
)

type Conversation struct {
	Id         string
	Channel    string
	Text       []string `gorm:"type:json;serializer:json"`
	StartedAt  string   `gorm:"type:string"`
	FinishedAt string   `gorm:"type:string"`
	Ver        int64
}

func NewConversation(ee stream.Events) ([]*Conversation, error) {
	var cc []*Conversation
	for _, e := range ee {
		cc = append(cc, &Conversation{Id: e.Root().Hash()})
	}

	return cc, nil
}

func (c *Conversation) ID() string { return c.Id }

func (c *Conversation) Version() int64 {
	return c.Ver
}

func (c *Conversation) Commit(event any, createdAt time.Time) error {
	//delay(time.Millisecond * 1500)
	switch e := event.(type) {
	case model.ThreadStarted:
		c.Channel, c.StartedAt = e.Channel, createdAt.String()

	case model.ThreadMessage:
		if strings.Contains(e.Text, "crush!") {
			return fmt.Errorf("oh no, it's crush message")
		}

		c.Text = append(c.Text, fmt.Sprintf("%s | %s |> %s",
			createdAt.Format(time.StampMilli),
			e.Participant,
			e.Text))

	case model.ThreadClosed:
		c.FinishedAt = createdAt.String()

	case model.ThreadLeft:
		//return fmt.Errorf("i will not accept it, do not want it")
	}

	c.Ver++

	return nil
}

func (c *Conversation) String() string {
	s := fmt.Sprintf("----- #%s.%s channel --------- %s ------------------\n",
		c.Channel, c.ID, c.StartedAt)
	for i := range c.Text {
		s += c.Text[i] + "\n"
	}

	s += fmt.Sprintf("----- #%s.%s channel --------- has %d messages ----------------------",
		c.Channel, c.ID, len(c.Text))

	return s
}

type Conversations struct {
	stream.Documents[*Conversation]
}

func NewConversations() *Conversations {
	s, err := storage[*Conversation](NewConversation)
	if err != nil {
		panic(err)
	}

	return &Conversations{s}
}

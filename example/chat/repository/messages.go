package repository

import (
	"fmt"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"github.com/sokool/stream/store/mysql"
	"os"
	"time"
)

type Messages struct {
	Id         string
	Channel    string
	Text       []string `gorm:"type:json;serializer:json"`
	StartedAt  string   `gorm:"type:string"`
	FinishedAt string   `gorm:"type:string"`
	Ver        int64
}

func NewMessage(m stream.Events) (*Messages, error) {
	id := m.Unique()
	if id.IsZero() {
		return nil, nil
	}

	return &Messages{Id: id.Hash()}, nil
}

func (c *Messages) ID() string { return c.Id }

func (c *Messages) Version() int64 {
	return c.Ver
}

func (c *Messages) Commit(event any, createdAt time.Time) error {
	switch e := event.(type) {
	case model.ThreadStarted:
		c.Channel, c.StartedAt = e.Channel, createdAt.String()

	case model.ThreadMessage:
		c.Text = append(c.Text, fmt.Sprintf("%s | %s |> %s",
			createdAt.Format(time.StampMilli),
			e.Participant,
			e.Text))

	case model.ThreadClosed:
		c.FinishedAt = createdAt.String()
	}

	c.Ver++

	return nil
}

func (c *Messages) String() string {
	s := fmt.Sprintf("----- #%s.%s channel --------- %s ------------------\n",
		c.Channel, c.ID, c.StartedAt)
	for i := range c.Text {
		s += c.Text[i] + "\n"
	}

	s += fmt.Sprintf("----- #%s.%s channel --------- has %d messages ----------------------",
		c.Channel, c.ID, len(c.Text))

	return s
}

type Messagesz stream.CRUD[*Messages]

func NewMessagez() Messagesz {
	if cdn := os.Getenv("MYSQL_EVENT_STORE"); cdn != "" {
		fmt.Println(cdn)
		c, err := mysql.NewConnection(cdn, &stream.Schemas{})
		if err != nil {
			panic(err)
		}

		m, err := mysql.NewTable[*Messages](c, NewMessage)
		if err != nil {
			panic(err)
		}

		return m
	}

	return stream.NewEntities[*Messages](NewMessage)
}

//type Channel struct {
//	ID   string
//	Name string
//}
//
//func (c Channel) Id() string {
//	return c.ID
//}
//
//func (c Channel) Append(event stream.Message) error {
//	return nil
//}

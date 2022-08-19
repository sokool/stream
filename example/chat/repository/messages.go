package repository

//
//import (
//	"fmt"
//	"time"
//
//	"github.com/livechat/flux-go/stream"
//)
//
//type Messages struct {
//	ID         string
//	Channel    string
//	Text       []string
//	StartedAt  *time.Time
//	FinishedAt *time.Time
//}
//
//func NewMessage(m stream.Message) *Messages {
//	if m.Stream.Name().String() != "Chat" {
//		return nil
//	}
//
//	return &Messages{ID: m.Stream.Value()}
//}
//
//func (c *Messages) Id() string { return c.ID }
//
//func (c *Messages) Append(event stream.Message) error {
//	switch e := event.Value.(type) {
//	case ThreadStarted:
//		c.ID, c.Channel, c.StartedAt = event.Stream.Value(), e.Channel, &event.CreatedAt
//
//	case Message:
//		c.Text = append(c.Text, fmt.Sprintf("%s | %s |> %s",
//			event.CreatedAt.Format(time.StampMilli),
//			e.Participant,
//			e.Text))
//
//	case ThreadClosed:
//		c.FinishedAt = &event.CreatedAt
//	}
//
//	return nil
//}
//
//func (c *Messages) String() string {
//	s := fmt.Sprintf("----- #%s.%s channel --------- %s ------------------\n",
//		c.Channel, c.ID, c.StartedAt.Format(time.StampMilli))
//	for i := range c.Text {
//		s += c.Text[i] + "\n"
//	}
//
//	s += fmt.Sprintf("----- #%s.%s channel --------- has %d messages ----------------------",
//		c.Channel, c.ID, len(c.Text))
//
//	return s
//}
//
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

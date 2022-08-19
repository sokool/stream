package repository

//
//import (
//	"fmt"
//	"time"
//
//	"github.com/livechat/flux-go/stream"
//)
//
//type Member struct {
//	Name   Person
//	Avatar string
//	MutedDue string
//	JoinedAt *time.Time
//	LeftAt   *time.Time
//}
//
//func (a *Member) Id() string {
//	return string(a.Name)
//}
//
//func (a *Member) Append(event stream.Message) error {
//	switch e := event.Value.(type) {
//	case ThreadJoined:
//		a.Name, a.JoinedAt = e.Participant, &event.CreatedAt
//
//	case ThreadLeft:
//		a.Name, a.LeftAt = e.Participant, &event.CreatedAt
//
//	case ThreadKicked:
//		a.Name, a.LeftAt = e.Participant, &event.CreatedAt
//
//	case ThreadMuted:
//		a.Name, a.MutedDue = e.Participant, e.Reason
//
//	default:
//		return nil
//	}
//
//	return nil
//}
//
//func (a *Member) String() string {
//	return fmt.Sprintf("%s |> %s", a.JoinedAt.Format(time.StampMilli), a.Name)
//}

package repository

import (
	"fmt"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"time"
)

type Member struct {
	Id       string
	Avatar   string
	MutedDue string
	JoinedAt time.Time `gorm:"type:string;serializer:json"`
	LeftAt   time.Time `gorm:"type:string;serializer:json"`
	Seq      int64
}

func NewMember(se stream.Events) ([]*Member, error) {
	var mm []*Member
	for i := range se {
		var id string
		switch e := se[i].Body().(type) {
		case model.ThreadJoined:
			id = string(e.Participant)

		case model.ThreadLeft:
			id = string(e.Participant)

		case model.ThreadKicked:
			id = string(e.Participant)

		case model.ThreadMuted:
			id = string(e.Participant)
		default:
			continue
		}
		mm = append(mm, &Member{Id: id})
	}
	return nil, nil
}

func (a *Member) ID() string {
	return string(a.Id)
}

func (a *Member) Version() int64 {
	return a.Seq
}

func (a *Member) Commit(event any, createdAt time.Time) error {
	switch e := event.(type) {
	case model.ThreadJoined:
		a.Id, a.JoinedAt = string(e.Participant), createdAt

	case model.ThreadLeft:
		a.Id, a.LeftAt = string(e.Participant), createdAt

	case model.ThreadKicked:
		a.Id, a.LeftAt = string(e.Participant), createdAt

	case model.ThreadMuted:
		a.Id, a.MutedDue = string(e.Participant), e.Reason
	}

	a.Seq++
	return nil
}

func (a *Member) String() string {
	if a == nil {
		return ""
	}
	return fmt.Sprintf("%s |> %s", a.JoinedAt.Format(time.StampMilli), a.Id)
}

type Members struct {
	stream.Documents[*Member]
}

func NewMembers() *Members {
	s, err := storage[*Member](NewMember)
	if err != nil {
		panic(err)
	}

	return &Members{s}
}

func (m *Members) Build(events <-chan stream.Events) error {
	//TODO implement me
	panic("implement me")
}

func (m *Members) Recent() ([]*Member, error) {
	qy := []byte(`SELECT * FROM member WHERE seq > 400`)
	return m.Documents.Read(qy)
}

func (m *Members) Name(n string) (*Member, error) {
	q := fmt.Sprintf(`SELECT id FROM member WHERE id = "%s"`, n)
	o, err := m.Documents.Read([]byte(q))
	if err != nil || len(o) == 0 {
		return nil, err
	}

	return o[0], nil
}

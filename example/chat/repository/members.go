package repository

import (
	"fmt"
	"time"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

type Members = stream.Projections[*Member]

func NewMembers(se *stream.Engine) (*Members, error) {
	var m Members
	var err error
	if m.Store, err = storage[*Member](NewMember); err != nil {
		return nil, err
	}
	return &m, m.Compose(se)
}

type Member struct {
	Id       threads.MemberID
	Avatar   string
	MutedDue string
	JoinedAt time.Time `gorm:"type:string;serializer:json"`
	LeftAt   time.Time `gorm:"type:string;serializer:json"`
	Seq      int64
}

func NewMember(se stream.Events) (*Member, error) {
	for i := range se {
		switch e := se[i].Body().(type) {
		case threads.ThreadJoined:
			return &Member{Id: e.Participant}, nil

		case threads.ThreadLeft:
			return &Member{Id: e.Participant}, nil

		case threads.ThreadKicked:
			return &Member{Id: e.Participant}, nil

		case threads.ThreadMuted:
			return &Member{Id: e.Participant}, nil
		}
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
	case threads.ThreadJoined:
		a.Id, a.JoinedAt = e.Participant, createdAt

	case threads.ThreadLeft:
		a.Id, a.LeftAt = e.Participant, createdAt

	case threads.ThreadKicked:
		a.Id, a.LeftAt = e.Participant, createdAt

	case threads.ThreadMuted:
		a.Id, a.MutedDue = e.Participant, e.Reason

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

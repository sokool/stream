package repository

import (
	"fmt"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"time"
)

type Member struct {
	Id       model.Member
	Avatar   string
	MutedDue string
	JoinedAt time.Time `gorm:"type:string;serializer:json"`
	LeftAt   time.Time `gorm:"type:string;serializer:json"`
	Seq      int64
}

func NewMember(se stream.Events) (*Member, error) {
	for i := range se {
		switch e := se[i].Body().(type) {
		case model.ThreadJoined:
			return &Member{Id: e.Participant}, nil

		case model.ThreadLeft:
			return &Member{Id: e.Participant}, nil

		case model.ThreadKicked:
			return &Member{Id: e.Participant}, nil

		case model.ThreadMuted:
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
	case model.ThreadJoined:
		a.Id, a.JoinedAt = e.Participant, createdAt

	case model.ThreadLeft:
		a.Id, a.LeftAt = e.Participant, createdAt

	case model.ThreadKicked:
		a.Id, a.LeftAt = e.Participant, createdAt

	case model.ThreadMuted:
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

type Members struct {
	*stream.Projection[*Member]
}

func NewMembers() *Members {
	s, err := storage[*Member](NewMember)
	if err != nil {
		panic(err)
	}

	return &Members{
		Projection: &stream.Projection[*Member]{
			Store: s,
		},
	}
}

package stream

import (
	"time"
)

type Root interface {
	//Entity
	ID() string
	Committer
	Uncommitted(clear bool) []event
}

type Committer interface {
	Commit(e event, createdAt time.Time) error //  todo not able to deny it (remove error)
}

type Command[R Root] func(R) error

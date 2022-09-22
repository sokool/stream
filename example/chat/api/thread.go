package api

import (
	"time"
)

type Thread struct{}

func (t *Thread) ID() string {
	panic("implement me")
}

func (t *Thread) Version() int64 {
	panic("implement me")
}

func (t *Thread) Commit(event any, createdAt time.Time) error {
	panic("implement me")
}

func (t *Thread) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (t *Thread) UnmarshalJSON(bytes []byte) error {
	panic("implement me")
}

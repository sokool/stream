package stream

import (
	"context"
	"time"
)

type Events interface {
	Stream(ID) ReadWriterAt
	Read(Query) Reader
	Write(m []Message) (n int, err error)
}

type Query struct {
	Stream     ID
	From, To   time.Time
	Descending bool
	Shutdown   context.Context
}

type events struct {
}

func NewEvents() Events {
	return &events{}
}

func (e *events) Stream(id ID) ReadWriterAt {
	//TODO implement me
	panic("implement me")
}

func (e *events) Read(query Query) Reader {
	//TODO implement me
	panic("implement me")
}

func (e *events) Write(m []Message) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

//func (e *events) ReadAll(desc bool) Reader {
//	//TODO implement me
//	panic("implement me")
//}

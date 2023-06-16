package stream_test

import (
	"testing"

	"github.com/icrowley/fake"
	. "github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

func TestNewAggregate(t *testing.T) {
	type scenario[R Root] struct {
		description string
		root        NewRoot[R]
		events      []any
		err         bool
	}

	cases := []scenario[*threads.Thread]{
		{"no root func fails", nil, threads.Events, true},
		{"no events fails", threads.New, nil, true},
		{"with root func and events is ok", threads.New, threads.Events, false},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			id, _ := NewID[threads.Thread](fake.CharactersN(8))
			_, err := NewAggregate(id, c.root, c.events)
			if c.err && err == nil {
				t.Fatalf("error expected")
			}
			if !c.err && err != nil {
				t.Fatalf("no error expected, got %v", err)
			}
		})
	}
}

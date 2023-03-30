package stream_test

import (
	"testing"

	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/threads"
)

func TestNewSequence(t *testing.T) {
	type scenario struct {
		description string
		id          string
		number      int64

		uuid   string
		string string
		err    bool
	}

	cases := []scenario{
		{
			description: "test id with no number",
			id:          "test",
			uuid:        "7f549bf5-988b-5826-9d08-59beec898623",
			string:      "4be0643f:Thread",
		},
		{
			description: "1hag8gh id with number 2",
			id:          "1hag8gh",
			number:      2,
			uuid:        "3e0d7469-3017-552c-a6fe-b8c1242523a3",
			string:      "fcee15e0:Thread#2",
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			s, err := stream.NewSequence[*threads.Thread](c.id, c.number)
			if c.number < 0 {
				s, err = stream.NewSequence[*threads.Thread](c.id)
				c.number = 0
			}
			if c.err && err == nil {
				t.Fatalf("error expected")
			}
			if !c.err && err != nil {
				t.Fatalf("no error expected, got %v", err)
			}
			if !s.UUID().Is(c.uuid) {
				t.Fatalf("expected %s uuid, got %s", c.uuid, s.UUID())
			}
			if s.Number() != c.number {
				t.Fatalf("number bleh")
			}
			if !s.Is(c.string) {
				t.Fatalf("expected %s string, got %s", c.string, s)
			}
		})
	}
}

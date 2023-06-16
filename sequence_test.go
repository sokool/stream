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

		// expected
		hash   string
		string string
		err    bool
	}

	cases := []scenario{
		{
			description: "test id with no number",
			id:          "test",
			hash:        "9bf8a226-2e21-5d3d-9c77-82b83c984458",
			string:      "4be0643f.Thread",
		},
		{
			description: "1hag8gh id with number 2",
			id:          "1hag8gh",
			number:      2,
			hash:        "3099d54f-011e-5a6f-9a9c-d2b0f0cfbaae",
			string:      "fcee15e0.Thread#2",
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
			if !s.Hash().Is(c.hash) {
				t.Fatalf("expected %s uuid, got %s", c.hash, s.Hash())
			}
			if s.Number() != c.number {
				t.Fatalf("number bleh")
			}
			if s.String() != c.string {
				t.Fatalf("expected %s string, got %s", c.string, s)
			}
		})
	}
}

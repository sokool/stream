package stream_test

import (
	"fmt"
	"testing"

	"github.com/sokool/stream"
)

func TestNewSequence(t *testing.T) {
	var err error

	if _, err = stream.NewSequence[string](""); err == nil {
		t.Fatal("error expected")
	}

	s, err := stream.NewSequence[string]("test")
	if err != nil {
		t.Fatal("error not expected")
	}

	fmt.Println(s, s.Next())
	s = s.Next()
	fmt.Println(s.ID())
}

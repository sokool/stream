package stream_test

import (
	"github.com/sokool/stream"
	"testing"
)

func TestNewType(t *testing.T) {
	st, err := stream.NewType("hi")
	if err != nil {
		t.Fatal(err)
	}

	if st.Hash() != "2025b67b-d5a7-53e6-960a-4c601c0f9455" {
		t.Fatal()
	}
}

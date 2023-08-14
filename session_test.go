package stream_test

import (
	"fmt"
	"testing"

	"github.com/sokool/stream"
)

func TestRole_String(t *testing.T) {
	r := stream.Role{
		ID:   "7sgs",
		Name: "customer",
	}

	r = r.Resource("7157", "quote", "")
	fmt.Println(r)
}

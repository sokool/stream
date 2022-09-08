package stream_test

import (
	"fmt"
	"github.com/sokool/stream"
	"github.com/sokool/stream/example/chat/model"
	"testing"
)

func Test(t *testing.T) {
	c, _ := model.NewThread("48")
	n, _ := stream.NewRootID(c)

	e, err := stream.NewEvent(n, model.ThreadStarted{}, 1)
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Println(e.ID())
	//fmt.Println(e.Sequence())
	//fmt.Println(e.Type())
	//fmt.Println(e.Body())
	//fmt.Println(e.Namespace())
	fmt.Println(e.GoString())

}

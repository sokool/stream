package mysql_test

import (
	"fmt"
	"github.com/sokool/stream"
	"github.com/sokool/stream/store/mysql"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	s, err := mysql.NewEventsStore(os.Getenv("MYSQL_EVENT_STORE"), nil)
	if err != nil {
		t.Fatal(err)
	}

	//id, _ := stream.ParseRootID("73HdaUj.Thread")
	q := stream.Query{
		//ID: id.ID(),
		//Root: id.Type(),
		//Events: []stream.Type{"ThreadStarted", "ThreadJoined", "ThreadMessage"},
		//FromSequence: 3,
		//Text:         "hi there",
		NewestFirst: true,
	}

	e := make(stream.Events, 6)
	//x, _ := stream.ParseRootID("51068.Application")
	fmt.Println(s.Reader(q).Read(e))
	//e = append(e, stream.Event[any]{
	//	ID:          "a",
	//	Type:        "Doa",
	//	Root:        x,
	//	Sequence:    1,
	//	Body:        nil,
	//	Meta:        stream.Meta{},
	//	Correlation: "",
	//	Causation:   "",
	//	CreatedAt:   time.Now(),
	//	Author:      "",
	//})
	fmt.Println(e.Unique().IsZero())
	fmt.Println(e.String())
	//s.ReadWriter(id).ReadAt(make(stream.Events, 5), 2)
	//fmt.Println(s.Reader(q).Read(make(stream.Events, 10)))
	//mysql.NewEventsReader(s.Connection, q).Read(make(stream.Events, 5))

}

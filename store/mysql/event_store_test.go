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

	//id, _ := stream.ParseRootID("analytics.Application")
	var q stream.Query

	//q.Root.ID, q.Root.Type, q.Root.Events = id.ID(), id.Type(), []stream.Type{"Registered", "Published"}

	q.Root.Events = []stream.Type{"Collected"}
	fmt.Println(s.Reader(q).Read(make(stream.Events, 10)))
	//mysql.NewEventsReader(s.Connection, q).Read(make(stream.Events, 5))

	//fmt.Printf("%#v\n", e)
}

package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/repository"
	"testing"
	"time"
)

func TestDocuments_Load(t *testing.T) {
	threads := repository.NewChats()
	m1 := repository.Member{Id: "Albert", Avatar: "elo.gif", Seq: 1, JoinedAt: time.Now()}
	m2 := repository.Member{Id: "Greg", Avatar: "greg.jpg", Seq: 1, JoinedAt: time.Now()}

	if err := threads.Members.Store.Update(&m1); err != nil {
		t.Fatal(err)
	}

	if err := threads.Members.Store.Update(&m2); err != nil {
		t.Fatal(err)
	}

	mm := make([]*repository.Member, 2)
	if err := threads.Members.Store.Read(mm, nil); err != nil {
		t.Fatal(err)
	}

	if len(mm) != 2 {
		t.Fatal()
	}

	m3 := repository.Member{Id: "Greg"}

	fmt.Println(threads.Members.Store.One(&m3))
	fmt.Println(m3)
	fmt.Println(mm)
	//fmt.Println(c[0].String())

}

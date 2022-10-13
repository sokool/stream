package stream_test

import (
	"fmt"
	"github.com/sokool/stream/example/chat/repository"
	"testing"
	"time"
)

func TestDocuments_Load(t *testing.T) {
	chats := repository.NewChats()
	m1 := repository.Member{Id: "Albert", Avatar: "elo.gif", Seq: 1, JoinedAt: time.Now()}
	m2 := repository.Member{Id: "Greg", Avatar: "greg.jpg", Seq: 1, JoinedAt: time.Now()}

	if err := chats.Members.Documents.Update(&m1); err != nil {
		t.Fatal(err)
	}

	if err := chats.Members.Documents.Update(&m2); err != nil {
		t.Fatal(err)
	}

	mm, err := chats.Members.Recent()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(chats.Members.Name("Greg"))
	//r
	//mm := make([]*repository.Member, 2)
	//if err := chats.Members.Documents.Read(mm, nil); err != nil {
	//	t.Fatal(err)
	//}
	//
	//if len(mm) != 2 {
	//	t.Fatal()
	//}

	m3 := repository.Member{Id: "Greg"}

	//fmt.Println(chats.Members.Documents.)
	fmt.Println(m3)
	for i := range mm {
		fmt.Println(mm[i].String())
	}

	//fmt.Println(c[0].String())

}

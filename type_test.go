package stream_test

import (
	"testing"

	"github.com/sokool/stream"
)

func TestNewType(t *testing.T) {
	type foo struct{}
	if s, err := stream.NewType[foo](); s.String() != "Foo" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType(foo{}); s.String() != "Foo" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType(&foo{}); s.String() != "Foo" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType[*foo](); s.String() != "Foo" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType(any(foo{})); s.String() != "Foo" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType[string](); s.String() != "String" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType("hi"); s.String() != "Hi" || err != nil {
		t.Fatal(err, s)
	}
	if s, err := stream.NewType("   nIcE    fORMat   "); s.String() != "NIcEfORMat" || err != nil {
		t.Fatal(err, s)
	}
	if _, err := stream.NewType(""); err == nil {
		t.Fatalf("error expected")
	}
}

func TestType_NewID(t *testing.T) {
	type lol struct{}
	s, err := stream.NewType[lol]()
	if err != nil {
		t.Fatal(err)
	}
	d, err := s.NewID("hi")
	if err != nil {
		t.Fatal(err)
	}
	if d.String() != "297f69d4.Lol" {
		t.Fatal()
	}
}

func TestType_Hash(t *testing.T) {
	type scenario struct {
		description string
		typ         stream.Type
	}
	cases := []scenario{
		{"scenario string", stream.MustType("scenario")},
		{"`   S cen  ar i o   ` string", stream.MustType("   S cen  ar i o   ")},
		{"scenario generic", stream.MustType[scenario]()},
		{"scenario struct", stream.MustType(scenario{})},
	}

	hash := "864f8517-abf8-578f-a294-82b3c51a54f6"
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if s := c.typ.Hash().String(); s != hash {
				t.Fatalf("expected `%s`, got `%s`", hash, s)
			}
		})
	}
}

package stream_test

import "testing"

func Test(t *testing.T) {
	type scenario struct {
		description string
	}

	cases := []scenario{
		{"my first subtest"},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			t.Fatalf("expected:%v got:%v", "", "")
		})
	}
}

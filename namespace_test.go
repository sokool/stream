package stream

import (
	"testing"
)

func TestNewNamespace(t *testing.T) {
	type (
		case_ struct {
			description string
			have        string
			want        string
		}
	)

	scenario := []case_{
		{description: "id.type is correct", have: "s73ava.Payment", want: "s73ava.Payment"},
	}

	for _, c := range scenario {
		t.Run(c.description, func(t *testing.T) {
			n, err := NewStringNamespace(c.have)
			if err != nil {
				t.Fatal(err)
			}

			if s := n.String(); s != c.want {
				t.Fatalf("expected:%v got:%v", c.want, s)
			}
		})
	}
}

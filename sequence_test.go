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
		{description: "type.id is correct", have: "typ.124", want: "124.Typ"},
		{description: "id|type is correct", have: "124|typ", want: "124.Typ"},
	}

	for _, c := range scenario {
		t.Run(c.description, func(t *testing.T) {
			n, err := ParseNamespace(c.have)
			if err != nil {
				t.Fatal(err)
			}

			if s := n.String(); s != c.want {
				t.Fatalf("expected:%v got:%v", c.want, s)
			}
		})
	}
}

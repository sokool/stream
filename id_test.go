package stream

import (
	"testing"
)

func TestFromString(t *testing.T) {
	type (
		have struct {
			stream string
		}

		want struct {
			id, typ string
		}

		case_ struct {
			description string
			have
			want
		}
	)

	scenario := []case_{
		{"type.id is correct", have{"name.124"}, want{"124", "Name"}},
		{"id|type is correct", have{"124|name"}, want{"124", "Name"}},
	}

	check := func(in have) (out want) {
		n, err := ParseID(in.stream)
		if err != nil {
			t.Fatal(err)
		}
		out.id, out.typ = n.Value(), n.Type().String()
		return
	}

	for _, c := range scenario {
		t.Run(c.description, func(t *testing.T) {
			if out := check(c.have); out != c.want {
				t.Fatalf("expected:%v got:%v", c.want, out)
			}
		})
	}
}

package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Sequence todo
type Sequence struct {
	id     ID
	number int64
}

func NewSequence[T any](id string) (Sequence, error) {
	var rid Sequence
	var err error

	if rid.id, err = NewID[T](id); err != nil {
		return rid, Err("invalid id string %w", err)
	}

	return rid, nil
}

func MustSequence[T any](id string) Sequence {
	rid, err := NewSequence[T](id)
	if err != nil {
		panic(err)
	}
	return rid
}

func ParseSequence(s string) (n Sequence, err error) {
	var p = strings.Split(s, ".")
	var l = len(p)
	if l < 2 {
		return n, Err("wrong `%s` format, please use <id>.<type>.<number> ie `N8hY13fsd.Chat.43`", s)
	}

	if n.id, err = ParseID(s); err != nil {
		return n, Err("parse sequence failed %w", err)
	}

	if l == 3 {
		m, err := strconv.Atoi(p[2])
		if err != nil {
			return n, Err("parse sequence number failed %w", err)
		}

		n.number = int64(m)
	}

	return
}

func (s Sequence) ID() ID {
	return s.id
}

func (s Sequence) Type() Type {
	return s.id.Type()
}

func (s Sequence) Number() int64 {
	return s.number
}

func (s Sequence) Next() Sequence {
	s.number++
	return s
}

func (s Sequence) UUID() UUID {
	return NewUUID(s.String())
}

func (s Sequence) MarshalJSON() ([]byte, error) {
	return json.Marshal(view{
		"ID":     s.ID().Value(),
		"Type":   s.Type(),
		"Number": s.Number(),
	})
}

func (s *Sequence) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		return nil
	}

	var v view
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	r, err := ParseSequence(fmt.Sprintf("%s.%s", v["ID"], v["Type"]))
	if err != nil {
		return err
	}

	*s = r
	return nil
}

func (s Sequence) String() string {
	if s.IsEmpty() {
		return ""
	}
	if s.number == 0 {
		return s.id.String()
	}
	return fmt.Sprintf("%s#%d", s.id, s.number)
}

func (s Sequence) IsEmpty() bool {
	return s.id.IsEmpty()
}

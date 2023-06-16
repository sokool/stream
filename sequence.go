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

func NewSequence[T any](id string, n ...int64) (s Sequence, err error) {
	if s.id, err = NewID[T](id); err != nil {
		return s, Err("invalid id string %w", err)
	}
	if len(n) > 0 {
		s.number = n[0]
	}
	return s, nil
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

func (s Sequence) Number() int64 {
	return s.number
}

func (s Sequence) Next() Sequence {
	s.number++
	return s
}

func (s Sequence) Hash() UUID {
	return NewUUID(fmt.Sprintf("%s:%d", s.id.Hash(), s.number))
}

func (s Sequence) MarshalJSON() ([]byte, error) {

	return json.Marshal(view{
		"ID":     s.id.UUID(),
		"Type":   s.id.Type(),
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

func (s Sequence) Reset() {
	s.number = 0
}

func (s Sequence) IsEmpty() bool {
	return s.id.IsEmpty()
}

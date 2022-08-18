package stream

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

// Sequence ...
type Sequence struct {
	namespace Namespace

	// sequence number in a stream.
	//
	// When zero, then Message is not transactional, therefore Writer can save
	// message without acknowledge or even remain it in memory for a while for
	// later persistence. It is faster to write such Message but not guaranteed
	// that will be stored on external device such as database, queue or file.
	//
	// When not zero, then stream is considered as transactional - guaranteed
	// that will be stored in exact place in a stream. Each new Message in
	// a stream should be sequential - having logical order. Reader should
	// respect this rule, and throw ErrWrongSequence error when next
	// message sequence in a stream is not in logical order.
	number int64
}

func NewSequence(n Namespace, number ...int64) (Sequence, error) {
	if len(number) == 0 {
		number = append(number, 0)
	}

	if number[0] < 0 {
		return Sequence{}, Err("sequence number between 0 and <max-int64> is required")
	}

	return Sequence{n, number[0]}, nil
}

func (s Sequence) ID() string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(s.String())).String()
}

func (s Sequence) Namespace() Namespace {
	return s.namespace
}

func (s Sequence) Number() int64 {
	return s.number
}

func (s Sequence) Next() Sequence {
	s.number++
	return s
}

func (s Sequence) Set(number int64) Sequence {
	s.number = number
	return s
}

func (s Sequence) IsZero() bool {
	return s.namespace.IsZero()
}

func (s Sequence) String() string {
	if s.IsZero() {
		return ""
	}

	return fmt.Sprintf("%s#%d", s.namespace, s.number)

}

func (s Sequence) MarshalJSON() ([]byte, error) {
	if s.IsZero() {
		return null, nil
	}
	return json.Marshal(s.String())
}

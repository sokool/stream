package stream

import "encoding/json"

type ID string

func NewID(s string) (ID, error) {
	if s == "" {
		return "", Err("id is empty")
	}
	return ID(s), nil
}

func (i ID) MarshalJSON() ([]byte, error) {
	if i == "" {
		return null, nil
	}

	return json.Marshal(i.String())
}

func (i ID) String() string {
	return string(i)
}

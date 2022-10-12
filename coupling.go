package stream

// todo map[ProjectionType][]EventType
type Types map[Type]bool

func NewTypes(s ...string) (Types, error) {
	c := map[Type]bool{}
	for i := range s {
		t, err := NewType(s[i])
		if err != nil {
			return nil, err
		}

		c[t] = true
	}

	return c, nil
}

func (c Types) merge(d Types) Types {
	for n := range d {
		c[n] = true
	}
	return c
}

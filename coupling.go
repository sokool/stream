package stream

type Coupling []Type

func NewCoupling(s ...string) (Coupling, error) {
	var c Coupling
	for i := range s {
		t, err := NewType(s[i])
		if err != nil {
			return nil, err
		}
		c = append(c, t)
	}

	return c, nil
}

func NewStrongCoupling() (Coupling, error) {
	return NewCoupling("@strong")
}

func (c Coupling) Add(t ...Type) Coupling {
	return append(c, t...)
}

func (c Coupling) IsLoose(with ...Type) bool {
	if len(c) == 0 {
		return true
	}

	for i := range c {
		if c[i] == "@strong" {
			return false
		}

		for k := range with {
			if c[i] == with[k] {
				return false
			}
		}
	}

	return true
}

func (c Coupling) IsStrong(with ...Type) bool {

	return !c.IsLoose(with...)
}

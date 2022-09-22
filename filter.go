package stream

type Filter interface {
	// Filtrate todo -> Filtrate([]Message) ([]Message, error)
	Filtrate(*Event[any]) (bool, error)
}

type FilterFunc func(*Event[any]) (bool, error)

func (fn FilterFunc) Filtrate(m *Event[any]) (bool, error) { return fn(m) }

type Filters []Filter

func (s Filters) Filtrate(m *Event[any]) (bool, error) {
	for i := range s {
		if ok, err := s[i].Filtrate(m); !ok || err != nil {
			return ok, err
		}
	}

	return true, nil
}

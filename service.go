package stream

type RR interface {
	Register(*Domain) error
}

type Domain struct {
	store   EventStore
	schemas Schemas
	logger  Logger
	writers *multiWriter
}

func NewDomain(es EventStore) *Domain {
	s := Domain{
		schemas: make(Schemas, 0),
		store:   NewEventStore(),
		logger:  DefaultLogger,
		writers: &multiWriter{},
	}

	if es != nil {
		s.store = es
	}

	return &s
}

func (s *Domain) Register(r ...RR) error {
	for i := range r {
		if err := r[i].Register(s); err != nil {
			return err
		}
	}
	return nil
}

// MultiWriter creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
//
// Each write is written to each listed writer, one at a time.
// If a listed writer returns an error, that overall write operation
// stops and returns the error; it does not continue down the list.
func MultiWriter(w ...Writer) Writer {
	s := make([]Writer, 0, len(w))
	for i := range w {
		if mw, ok := w[i].(*multiWriter); ok {
			s = append(s, mw.list...)
		} else {
			s = append(s, w[i])
		}
	}
	return &multiWriter{s}
}

type multiWriter struct {
	list []Writer
}

func (m *multiWriter) Write(e Events) (n int, err error) {
	for _, w := range m.list {
		if n, err = w.Write(e); err != nil {
			return
		}
		if n != len(e) {
			err = ErrShortWrite
			return
		}
	}
	return len(e), nil
}

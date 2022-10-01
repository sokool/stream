package stream

import "os"

type Registerer interface {
	Register(*Domain) error
}

type Domain struct {
	store   EventStore
	logger  Logger
	writers *multiWriter
}

type Configuration struct {
	// Logger
	Logger func(Type) Printer

	// EventStore factory
	EventStore func(Printer) EventStore // todo func not needed
}

func NewDomain(c *Configuration) *Domain {
	s := Domain{
		store:   NewEventStore(),
		logger:  NewLogger(os.Stdout, "stream", true).WithTag,
		writers: &multiWriter{},
	}

	if c.Logger != nil {
		s.logger = c.Logger
	}
	if c.EventStore != nil {
		s.store = c.EventStore(s.logger("EventStore"))
	}

	return &s
}

func (s *Domain) Register(r ...Registerer) error {
	for i := range r {
		if err := r[i].Register(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *Domain) Run() {}

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

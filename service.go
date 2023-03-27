package stream

import (
	"os"
	"sync"
)

type Component interface {
	Compose(*Engine) error
}

type Configuration struct {
	Name Type
	// Logger
	Logger func(Type) Printer

	// EventStore factory
	EventStore func(Printer) EventStore // todo func not needed
}

type Engine struct {
	name    Type
	store   EventStore
	logger  Logger
	log     Printer
	writers map[Type]Writer
	mu      sync.RWMutex
}

func New(c *Configuration) *Engine {
	s := Engine{
		name:    c.Name,
		store:   MemoryEventStore,
		logger:  NewLogger(os.Stdout, "stream", true).WithTag,
		writers: map[Type]Writer{},
	}

	if c.Logger != nil {
		s.logger = c.Logger
	}
	if c.EventStore != nil {
		s.store = c.EventStore(s.logger("EventStore"))
	}

	s.log = s.logger(s.name)

	return &s
}

func (s *Engine) Compose(c ...Component) error {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	for i := range c {
		if err := c[i].Compose(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *Engine) Write(e Events) (n int, err error) {
	var swg sync.WaitGroup
	var che = make(chan error, len(s.writers))

	s.mu.Lock()
	defer s.mu.Unlock()

	for t := range s.writers {
		swg.Add(1)
		go func(t Type, w Writer, e Events) {
			ok := registry.isCoupled(t, e)
			if !ok {
				swg.Done()
			}

			if _, err = w.Write(e); err != nil {
				if ok {
					che <- err
				} else {
					s.log("ERR %s %s failed due `%s` error", e, t, err)
				}
			}

			if ok {
				swg.Done()
			}

		}(t, s.writers[t], e)
	}

	swg.Wait()
	close(che)

	return len(e), <-che
}

func (s *Engine) register(w Writer, t Type) error {
	if _, ok := s.writers[t]; ok {
		return Err("%s already registered", t)
	}
	s.writers[t] = w

	return nil
}

func (s *Engine) Run() {}

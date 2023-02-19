package stream

import (
	"reflect"
	"sync"
	"time"
)

type Document interface {
	Entity
	Committer
}

type Projections[D Document] struct {
	// Name is unique identifier of events handler
	Name Type

	// Description
	Description string

	// OnEvents
	// or Entities must be set
	OnEvents func(Events) (D, error)

	// OnFilter
	OnFilter Filter

	// OnEvent
	//OnEvent AppenderFunc

	// OnBuild
	//OnBuild Receiver

	// BuildOnStart
	BuildOnStart bool

	// BuildLogRefresh
	BuildLogRefresh time.Duration

	// Logger
	Log Printer

	Store Entities[D]

	mu      sync.Mutex
	blocked error
}

func (p *Projections[D]) init() error {
	if p.Name.IsZero() {
		var v D
		var err error
		if p.Name, err = NewType(v); err != nil {
			return err
		}
	}

	if p.OnEvents == nil && p.Store == nil {
		return Err("%s projection requires Document CRUD implementation or OnEvents func", p.Name)
	}

	if p.Log == nil {
		p.Log = DefaultLogger(p.Name)
	}

	//p.log("initialized queue size: %d, delivery timeout: %s", p.EventsQueueSize, p.EventsDeliveryTimeout)
	return nil
}

func (p *Projections[D]) Write(e Events) (n int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err = p.init(); err != nil {
		return 0, err
	}

	if n = len(e); n == 0 {
		return
	}

	if p.blocked != nil {
		return 0, p.blocked // return ErrBlocker or somthing
	}

	defer func(t time.Time) {
		if err != nil {
			p.blocked = err
			return
		}
		p.log("DBG %s delivered in %s", e, time.Since(t))
	}(time.Now())

	return n, p.write(e)
}

func (p *Projections[D]) write(e Events) (err error) {
	//if e, err = e.Shrink(h.OnFilter); err != nil {
	//	return 0, err
	//}

	//if h.builder != nil && h.builder.IsRunning() {
	//	if err = h.builder.Append(m); err != nil {
	//		return n, err
	//	}
	//
	//	return n, nil
	//}
	//}

	var d D
	if p.OnEvents != nil {
		if d, err = p.OnEvents(e); err != nil {
			return err
		}
	}

	if p.Store != nil {
		d, err = p.Store.Create(e)
		if err != nil || reflect.ValueOf(d).IsNil() { //todo i do now how to check generic D is nil
			return err
		}

		if err = p.Store.One(d); err != nil && err != ErrDocumentNotFound {
			return err
		}

		for i := range e {
			if err = d.Commit(e[i].body, e[i].createdAt); err != nil {
				return err
			}
		}

		if err = p.Store.Update(d); err != nil {
			return err
		}
	}

	return nil
}

func (p *Projections[D]) log(m string, a ...interface{}) {
	if p.Log == nil {
		return
	}

	p.Log(m, a...)
}

func (p *Projections[D]) Compose(in *Service) error {
	if err := p.init(); err != nil {
		return err
	}

	if err := in.register(p, p.Name); err != nil {
		return err
	}
	p.Log("projection composed")
	return nil
}

//func (h *Projections[D]) query() Query {
//	if q, ok := h.OnFilter.(*Query); ok {
//		return *q
//	}
//
//	if q, ok := h.OnFilter.(Definitions); ok {
//		return q.Query()
//	}
//
//	return Query{}
//}

//type Handlers struct {
//	mu         sync.Mutex
//	store      EventStore
//	schemas    *Schemas
//	registered map[string]*Projections
//	log        Logger
//	w          Writer //todo it support old projections
//}
//
//func NewHandlers(e EventStore, s *Schemas, l Logger, w Writer) *Handlers {
//	return &Handlers{
//		store:      e,
//		schemas:    s,
//		log:        l,
//		w:          w,
//		registered: map[string]*Projections{},
//	}
//}
//
//func (r *Handlers) Set(m Events) (n int, err error) {
//	if len(m) == 0 {
//		return 0, nil
//	}
//
//	for _, h := range r.registered {
//		if n, err = h.Set(m); err != nil {
//			h.log("ERR %s delivery failed", MString(m))
//			continue
//		}
//	}
//
//	if r.w != nil {
//		go func(w Writer, ee Events) {
//			if _, failed := w.Set(ee); failed != nil {
//				r.log("handlers")("ERR %s", failed)
//			}
//		}(r.w, m)
//	}
//
//	return len(m), nil
//}
//
//func (r *Handlers) Register(h ...*Projections) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	for i := range h {
//		p := h[i]
//
//		if p.Type == "" {
//			return Err("stream handlers: name is required")
//		}
//
//		if _, found := r.registered[p.Type]; found {
//			return Err("stream handlers: %s already exists", p.Type)
//		}
//
//		if p.Log == nil {
//			p.Log = r.log(p.Type)
//		}
//
//		if p.schemas == nil {
//			p.schemas = r.schemas
//		}
//
//		r.registered[p.Type] = p
//
//		if p.BuildOnStart {
//			b := p.Builder()
//			if b == nil {
//				return Err("%s building not supported", p.Type)
//			}
//
//			if err := b.Start(r.store.Search(p.query())); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
//
//func (r *Handlers) Get(name string) *Projections {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	return r.registered[name]
//}

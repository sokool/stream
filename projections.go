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

type NewDocument[D Document] func(Events) (D, error)

type Projections[D Document] struct {
	mu sync.Mutex

	// typ is unique identifier of events handler
	typ Type

	// description
	//description string

	// onCreate
	// or Entities must be set
	onCreate NewDocument[D]

	// onFilter
	//onFilter Filter

	// OnEvent
	//OnEvent AppenderFunc

	// OnBuild
	//OnBuild Receiver

	// buildOnStart
	//buildOnStart bool

	// buildLogRefresh
	//buildLogRefresh time.Duration

	// Logger
	log Printer

	Store   Entities[D]
	blocked error
}

func NewProjections[D Document](nd NewDocument[D]) *Projections[D] {
	dt := MustType[D]()
	return &Projections[D]{
		typ:      dt,
		onCreate: nd,
		log:      DefaultLogger(dt),
		Store:    NewEntities[D](),
	}
}

func (p *Projections[D]) init() error {
	if p.onCreate == nil && p.Store == nil {
		return Err("%s projection requires Document CRUD implementation or OnEvents func", p.typ)
	}

	//p.log("initialized queue size: %d, delivery timeout: %s", p.EventsQueueSize, p.EventsDeliveryTimeout)
	return nil
}

func (p *Projections[D]) Write(e Events) (n int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

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

func (p *Projections[D]) WithLogger(l Logger) *Projections[D] {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log = l(p.typ)
	return p
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
	if p.onCreate != nil {
		if d, err = p.onCreate(e); err != nil {
			return err
		}
	}

	if p.Store != nil {
		d, err = p.onCreate(e)
		//d, err = p.Store.Create(e)
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

func (p *Projections[D]) Compose(in *Engine) error {
	if err := in.register(p, p.typ); err != nil {
		return err
	}

	p.WithLogger(in.logger)
	p.log("projection composed")
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

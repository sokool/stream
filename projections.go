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

type Projection[D Document] struct {
	// Name is unique identifier of events handler
	Name Type

	// Description
	Description string

	// OnEvents must be set or Entities interface
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

	Documents Entities[D]

	mu sync.Mutex

	//builder *Builder
	schemas Schemas
	time    time.Duration
	recent  *Event[any]
	count   int64
}

func (h *Projection[D]) Write(e Events) (n int, err error) {
	s := len(e)
	if s == 0 {
		return 0, nil
	}

	if h.schemas.IsCoupled(h.Name, e) {
		return h.write(e)
	}

	go func(m Events) {
		if _, failed := h.write(m); failed != nil {
			h.log("ERR %s failed due %s", m, failed)
		}
	}(e)

	return s, nil
}

func (h *Projection[D]) write(e Events) (_ int, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if e, err = e.Shrink(h.OnFilter); err != nil {
		return 0, err
	}

	//if h.builder != nil && h.builder.IsRunning() {
	//	if err = h.builder.Append(m); err != nil {
	//		return n, err
	//	}
	//
	//	return n, nil
	//}
	//}
	var d D
	if h.OnEvents != nil {
		if d, err = h.OnEvents(e); err != nil {
			return 0, err
		}
	}

	if h.Documents != nil {

		d, err = h.Documents.Create(e)
		if err != nil || reflect.ValueOf(d).IsNil() { //todo i do now how to assert generic D to nil
			return 0, err
		}

		if err = h.Documents.One(d); err != nil && err != ErrDocumentNotFound {
			return 0, err
		}

		for i := range e {
			if err = d.Commit(e[i].Body, e[i].CreatedAt); err != nil {
				return 0, err
			}
		}

		if err = h.Documents.Update(d); err != nil {
			return 0, err
		}

	}

	h.log("DBG %s delivered", e)
	return len(e), nil
}

func (h *Projection[D]) log(m string, a ...interface{}) {
	if h.Log == nil {
		return
	}
	h.Log(m, a...)
}

func (h *Projection[D]) Register(in *Domain) (err error) {
	if h.Name.IsZero() {
		var v D
		if h.Name, err = NewType(v); err != nil {
			return err
		}
	}

	if h.OnEvents == nil && h.Documents == nil {
		return Err("%s projection requires Document CRUD implementation or OnEvents func", h.Name)
	}

	if h.Log == nil {
		h.Log = in.logger(h.Name.String())
	}

	in.writers.list = append(in.writers.list, h)
	return nil
}

//func (h *Projection[D]) query() Query {
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
//	registered map[string]*Projection
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
//		registered: map[string]*Projection{},
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
//func (r *Handlers) Register(h ...*Projection) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	for i := range h {
//		p := h[i]
//
//		if p.Name == "" {
//			return Err("stream handlers: name is required")
//		}
//
//		if _, found := r.registered[p.Name]; found {
//			return Err("stream handlers: %s already exists", p.Name)
//		}
//
//		if p.Log == nil {
//			p.Log = r.log(p.Name)
//		}
//
//		if p.schemas == nil {
//			p.schemas = r.schemas
//		}
//
//		r.registered[p.Name] = p
//
//		if p.BuildOnStart {
//			b := p.Builder()
//			if b == nil {
//				return Err("%s building not supported", p.Name)
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
//func (r *Handlers) Get(name string) *Projection {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	return r.registered[name]
//}

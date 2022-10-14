package stream

import "sync"

type Handler struct {
	// Name is unique identifier of events handler
	Name Type

	// Description
	Description string

	// OnEvents
	OnEvents func(Events) error

	// OnFilter
	OnFilter Filter

	mu sync.Mutex
}

func (h *Handler) Write(e Events) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.OnFilter != nil {
		if ok, err := h.OnFilter.Filtrate(e); err != nil || !ok {
			return 0, err
		}
	}

	if err := h.OnEvents(e); err != nil {
		return 0, err
	}

	return e.Size(), nil
}

func (h *Handler) register(d *Domain) error {
	if h.Name.IsZero() {
		return Err("handler name required")
	}

	if h.OnEvents == nil {
		return Err("handler OnEvents required")
	}

	return d.register(h, h.Name)
}

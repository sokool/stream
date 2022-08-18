package stream

type Meta struct {
	Author      string
	Correlation ID
}

func NewMeta(c Context) (m Meta, err error) {
	var text = func(key string) string { s, _ := c.Value(key).(string); return s }

	m.Author = text("stream-author")

	if s := text("stream-correlation-id"); s != "" {
		if m.Correlation, err = NewID(s); err != nil {
			return
		}
	}

	return
}

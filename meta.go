package stream

type Meta struct {

	// Every Message has 3 ID's [ID, CorrelationID, CausationID]. When you are
	// responding to a Message (either a Thread or and Event) you copy the
	// CorrelationID of the Message you are responding to, to your new
	// CorrelationID. The CausationID of your Message is the ID of the
	// Message you are responding to.
	//
	// Greg Young
	// --> https://groups.google.com/d/msg/dddcqrs/qGYC6qZEqOI/LhQup9v7EwAJ
	Correlation, Causation ID

	// Author helps to check what person/device generate this Message.
	Author string
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

package stream

import "time"

type Message struct {
	ID
	Type

	// body TODO
	body []byte

	// meta TODO
	meta []byte

	// sequence number in a stream.
	//
	// When zero, then Message is not transactional, therefore Writer can save
	// message without acknowledge or even remain it in memory for a while for
	// later persistence. It is faster to write such Message but not guaranteed
	// that will be stored on external device such as database, queue or file.
	//
	// When not zero, then stream is considered as transactional - guaranteed
	// that will be stored in exact place in a stream. Each new Message in
	// a stream should be sequential - having logical order. Reader should
	// respect this rule, and throw ErrWrongSequence error when next
	// message sequence in a stream is not in logical order.
	sequence int64

	// Every Message has 3 ID's [ID, CorrelationID, CausationID]. When you are
	// responding to a Message (either a Command or and Event) you copy the
	// CorrelationID of the Message you are responding to, to your new
	// CorrelationID. The CausationID of your Message is the ID of the
	// Message you are responding to.
	//
	// Greg Young
	// --> https://groups.google.com/d/msg/dddcqrs/qGYC6qZEqOI/LhQup9v7EwAJ
	IDx, CorrelationID, CausationID string

	// CreatedAt
	CreatedAt time.Time

	// Author helps to check what person/device generate this Message.
	Author string

	// Payload TODO
	value interface{} `json:"-"`
}

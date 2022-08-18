package stream

// Reader
type Reader[E any] interface {
	Read(m []Message[E]) (n int, err error)
}

type ReadWriterCloser[E any] interface {
	Reader[E]
	Writer[E]
	Closer
}

// ReaderAt
type ReaderAt[E any] interface {
	ReadAt(m []Message[E], pos int64) (n int, err error)
}

type ReaderFrom[E any] interface {
	ReadFrom(r Reader[E]) (n int64, err error)
}

type Closer interface {
	Close(error) error
}

type ReadWriterAt[E any] interface {
	ReaderAt[E]
	WriterAt[E]
}

type Writer[E any] interface {
	Write(m []Message[E]) (n int, err error)
}

type WriteCloser[E any] interface {
	Writer[E]
	Closer
}

// WriterAt store messages in Entity starting from pos... todo
type WriterAt[E any] interface {
	WriteAt(m []Message[E], pos int64) (n int, err error)
}

type WriterTo[E any] interface {
	WriteTo(w Writer[E]) (n int64, err error)
}

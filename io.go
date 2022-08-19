package stream

// Reader
type Reader[E any] interface {
	Read(m []Event[E]) (n int, err error)
}

type ReadWriterCloser[E any] interface {
	Reader[E]
	Writer[E]
	Closer
}

// ReaderAt
type ReaderAt[E any] interface {
	ReadAt(m []Event[E], pos int64) (n int, err error)
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
	Write(m []Event[E]) (n int, err error)
}

type WriteCloser[E any] interface {
	Writer[E]
	Closer
}

// WriterAt store event in Entity starting from pos... todo
type WriterAt[E any] interface {
	WriteAt(m []Event[E], pos int64) (n int, err error)
}

type WriterTo[E any] interface {
	WriteTo(w Writer[E]) (n int64, err error)
}

package stream

// Reader
type Reader interface {
	Read(m []Event[any]) (n int, err error)
}

type ReadWriterCloser interface {
	Reader
	Writer
	Closer
}

// ReaderAt
type ReaderAt interface {
	ReadAt(m []Event[any], pos int64) (n int, err error)
}

type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

type Closer interface {
	Close(error) error
}

type ReadWriterAt interface {
	ReaderAt
	WriterAt
}

type Writer interface {
	Write(m []Event[any]) (n int, err error)
}

type WriteCloser[E any] interface {
	Writer
	Closer
}

// WriterAt store event in Entity starting from pos... todo
type WriterAt interface {
	WriteAt(m []Event[any], pos int64) (n int, err error)
}

type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

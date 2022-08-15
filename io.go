package stream

// Reader
type Reader interface {
	Read(m []Message) (n int, err error)
}

type ReadWriterCloser interface {
	Reader
	Writer
	Closer
}

// ReaderAt
type ReaderAt interface {
	ReadAt(m []Message, pos int64) (n int, err error)
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
	Write(m []Message) (n int, err error)
}

type WriteCloser interface {
	Writer
	Closer
}

// WriterAt store messages in Entity starting from pos... todo
type WriterAt interface {
	WriteAt(m []Message, pos int64) (n int, err error)
}

type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

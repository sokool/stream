package stream

// Reader
type Reader interface {
	Read(Events) (n int, err error)
}

type ReadWriterCloser interface {
	Reader
	Writer
	Closer
}

// ReaderAt
// todo replace pos with Sequence
type ReaderAt interface {
	ReadAt(e Events, pos int64) (n int, err error)
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
	Write(Events) (n int, err error)
}

type WriteCloser interface {
	Writer
	Closer
}

// WriterAt store event in Document starting from pos... todo
// todo replace pos with Sequence
type WriterAt interface {
	WriteAt(m Events, pos int64) (n int, err error)
}

type WriterTo interface {
	WriteTo(Writer) (n int64, err error)
}

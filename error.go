package stream

import (
	"github.com/sokool/errors"
)

var Err = errors.Errorf

var (
	// ErrEndOfStream is the error returned by Reader when no more input is available.
	// Reading functions should return ErrEndOfStream only to signal a graceful end of input.
	// If the ErrEndOfStream occurs unexpectedly in a structured data stream,
	// the appropriate error is either ErrWrongPosition or some other error related
	// to underlying Reader
	ErrEndOfStream = Err("stream: end of stream")

	ErrDocumentNotFound = Err("stream:document not found")
	// ErrWrongSequence error might be returned by Reader or Writer.
	// They can detect if Messages are in logical order, when
	// not then ErrWrongSequence should be returned.
	//
	// Messages are organized in group of stream...
	//ErrWrongSequence = Err("sequence problem in a stream")

	// ErrConcurrentWrite when Appender or Writer is running.
	ErrConcurrentWrite = Err("stream: concurrent write")
	//
	// ErrDuplicatedMessage
	//ErrDuplicatedMessage = Err("duplicated message in a stream")

	// ErrShortWrite means that a write accepted fewer Message number than requested
	// but failed to return an explicit error.
	ErrShortWrite = Err("stream: short write")

	//ErrPaused = Err("stream paused")

	//ErrCancelled = Err("stream cancelled")

	// ErrBuildInProgress
	//ErrBuildInProgress = Err("building in progress")
)

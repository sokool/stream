package stream

import "fmt"

var Err = func(format string, args ...any) error {
	return fmt.Errorf("stream:"+format, args...)
}

package stream

import (
	"bytes"
	"fmt"
	"io"
	"log/syslog"
	"os"
	"strings"
	"time"
)

type Printer = func(string, ...interface{})

type Logger = func(tag string) Printer

var DefaultLogger = func(tag string) Printer {
	return NewLogger(os.Stdout, tag, true).Print
}

type logger struct {
	w         io.Writer
	verbose   bool
	colors    bool
	timestamp bool
	tag       string
}

func NewLogger(w io.Writer, tag string, verbose bool) *logger {
	_, timestamp := w.(*os.File)
	if tag != "" {
		tag = fmt.Sprintf("%s", tag)
	}

	return &logger{
		w:         w,
		verbose:   verbose,
		colors:    w == os.Stdout,
		tag:       tag,
		timestamp: timestamp,
	}
}

func (l *logger) Write(p []byte) (n int, err error) {
	if s, ok := l.contains(p, "dbg:", "dbg", "[DEBUG]"); ok {
		l.Print("DBG " + s)
		return
	}

	if s, ok := l.contains(p, "err:", "err", "[ERROR]"); ok {
		l.Print("ERR " + s)
		return
	}

	if s, ok := l.contains(p, "inf:", "nfo:", "inf", "info"); ok {
		l.Print("INF " + s)
		return
	}

	l.Print("INF " + string(p))
	return
}

func (l *logger) Print(format string, a ...interface{}) {
	s := strings.Split(format, " ")
	typ := strings.ToUpper(s[0])

	if typ != "INF" && typ != "DBG" && typ != "ERR" {
		format = "INF " + format
		s[0] = "INF"
		typ = "INF"
	}

	format = strings.TrimSpace(strings.Replace(format, s[0], "", 1))

	if len(a) >= 1 {
		if _, ok := a[0].(error); ok {
			typ = "ERR"
		}
	}

	if typ == "DBG" && !l.verbose {
		return
	}

	// syslog support
	if w, ok := l.w.(*syslog.Writer); ok {
		m := fmt.Sprintf("%s %s %s", typ, l.tag, fmt.Sprintf(format, a...))

		switch typ {
		case "INF":
			w.Info(m)
		case "ERR":
			w.Err(m)
		case "DBG":
			w.Debug(m)
		}

		return
	}

	color := "%s"
	if l.colors {
		switch typ {
		case "INF":
			color = "\x1b[32;1m%s\x1b[0m" // green

		case "ERR":
			color = "\x1b[31;1m%s\x1b[0m" // red

		case "DBG":
			color = "\x1b[33;1m%s\x1b[0m" // yellow
		}
	}

	format = fmt.Sprintf(format, a...)
	format = strings.TrimSuffix(format, "\n")
	typ = fmt.Sprintf(color, typ)
	x := l.tag
	if l.tag != "" {
		x = fmt.Sprintf("[\x1b[36;1m%s\x1b[0m] ", l.tag)
	}
	n := time.Now().Format("2006/01/02 15:04:05.000")
	l.w.Write([]byte(fmt.Sprintf("%s [%s] %s%s\n", n, typ, x, format)))
}

func (l *logger) WithTag(name string) Printer { return NewLogger(l.w, name, l.verbose).Print }

func (l *logger) contains(b []byte, of ...string) (string, bool) {
	for i := range of {
		if k := bytes.Index(bytes.ToLower(b), []byte(of[i])); k == 0 {
			return string(b[len(of[i]):]), true
		}
	}

	return "0", false
}

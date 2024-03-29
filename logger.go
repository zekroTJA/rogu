package rogu

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/zekrotja/rogu/level"
)

// Logger is used to create log events and
// to specify writers which are used to
// write commited events.
type Logger interface {
	Closer

	slog.Handler

	AddWriter(w Writer) Logger
	Copy() *logger
	Debug() *Event
	Error() *Event
	Fatal() *Event
	Info() *Event
	Panic() *Event
	SetCaller(enable bool) Logger
	SetLevel(lvl level.Level) Logger
	SetWriter(w Writer) Logger
	Tagged(tag string) Logger
	Trace() *Event
	Warn() *Event
	WithLevel(lvl level.Level) *Event
}

type eventWriter interface {
	write(e *Event, msg string) error
}

type logger struct {
	w      Writer
	lvl    level.Level
	caller bool
}

var _ Logger = (*logger)(nil)

// NewLogger returns a new instance of Logger
// with the passed writers.
//
// If no writer is specified or set via `SetWriter`,
// the logger will never output anything.
func NewLogger(writer ...Writer) Logger {
	l := &logger{}

	if len(writer) == 1 {
		l.w = writer[0]
	} else {
		l.w = MultiWriter(writer)
	}

	l.SetLevel(level.Info)

	return l
}

// SetWriter sets the specified writer to
// the logger.
func (t *logger) SetWriter(w Writer) Logger {
	t.w = w
	return t
}

// AddWriter adds another writer to the logger.
func (t *logger) AddWriter(w Writer) Logger {
	if t.w == nil {
		t.w = w
	} else if mw, ok := t.w.(MultiWriter); ok {
		t.w = MultiWriter(append(mw, w))
	} else {
		t.w = MultiWriter{t.w, w}
	}
	return t
}

// SetLevel sets the minum log leven which
// will be written.
func (t *logger) SetLevel(lvl level.Level) Logger {
	t.lvl = lvl
	return t
}

// SetCaller enabled or disables attaching the
// caller file and line to the event.
func (t *logger) SetCaller(enable bool) Logger {
	t.caller = enable
	return t
}

// Copy creates and returns a copy of the Logger.
func (t *logger) Copy() *logger {
	n := *t
	return &n
}

// Tagged returns a new logger which references
// the origin logger but attaches the given
// tag to every created Entry. Changes made
// to the underlying logger will be projected
// to the created logger.
func (t *logger) Tagged(tag string) Logger {
	n := &taggedLogger{
		logger: t,
		tag:    tag,
	}
	n.tag = tag
	return n
}

// Close closes the set writers or all writers that
// are added to the logger and which are closable.
func (t *logger) Close() error {
	if closer, ok := t.w.(Closer); ok {
		return closer.Close()
	}
	return nil
}

// Trace creates a new log Event with level trace.
func (t *logger) Trace() *Event {
	return t.newEvent(level.Trace)
}

// Trace creates a new log Event with level debug.
func (t *logger) Debug() *Event {
	return t.newEvent(level.Debug)
}

// Trace creates a new log Event with info.
func (t *logger) Info() *Event {
	return t.newEvent(level.Info)
}

// Trace creates a new log Event with level warn.
func (t *logger) Warn() *Event {
	return t.newEvent(level.Warn)
}

// Trace creates a new log Event with level error.
func (t *logger) Error() *Event {
	return t.newEvent(level.Error)
}

// Trace creates a new log Event with level fatal.
//
// When commited, the programm will exit with exit
// code 1.
func (t *logger) Fatal() *Event {
	return t.newEvent(level.Fatal)
}

// Trace creates a new log Event with level panic.
//
// When commited, the program will panic at the
// called point.
func (t *logger) Panic() *Event {
	return t.newEvent(level.Panic)
}

// WithLevel returns a new log Event with the given level.
func (t *logger) WithLevel(lvl level.Level) *Event {
	return t.newEvent(lvl)
}

func (t *logger) newEvent(lvl level.Level) *Event {
	e := newEvent(t, lvl)
	if t.caller {
		e.Caller()
	}
	return e
}

func (t *logger) write(e *Event, msg string) error {
	if e.lvl == level.Fatal {
		defer os.Exit(1)
	} else if e.lvl == level.Panic {
		defer panic(msg)
	}

	if e.lvl > t.lvl {
		return nil
	}

	if t.w == nil {
		return nil
	}

	var (
		file string
		line int
	)
	if e.caller {
		_, file, line, _ = runtime.Caller(2)
	}

	return t.w.Write(
		e.lvl,
		e.fields,
		e.tag,
		e.err,
		e.errFormat,
		file,
		line,
		msg,
	)
}

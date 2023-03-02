package rogu

import (
	"github.com/zekrotja/rogu/level"
)

type taggedLogger struct {
	*logger
	tag string
}

var _ Logger = (*taggedLogger)(nil)

// Trace creates a new log Event with level trace.
func (t *taggedLogger) Trace() *Event {
	return t.newEvent(level.Trace).Tag(t.tag)
}

// Trace creates a new log Event with level debug.
func (t *taggedLogger) Debug() *Event {
	return t.newEvent(level.Debug).Tag(t.tag)
}

// Trace creates a new log Event with info.
func (t *taggedLogger) Info() *Event {
	return t.newEvent(level.Info).Tag(t.tag)
}

// Trace creates a new log Event with level warn.
func (t *taggedLogger) Warn() *Event {
	return t.newEvent(level.Warn).Tag(t.tag)
}

// Trace creates a new log Event with level error.
func (t *taggedLogger) Error() *Event {
	return t.newEvent(level.Error).Tag(t.tag)
}

// Trace creates a new log Event with level fatal.
//
// When commited, the programm will exit with exit
// code 1.
func (t *taggedLogger) Fatal() *Event {
	return t.newEvent(level.Fatal).Tag(t.tag)
}

// Trace creates a new log Event with level panic.
//
// When commited, the program will panic at the
// called point.
func (t *taggedLogger) Panic() *Event {
	return t.newEvent(level.Panic).Tag(t.tag)
}

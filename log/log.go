package log

import (
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/level"
)

var defaultLogger = rogu.NewLogger(rogu.NewPrettyWriter())

// SetWriter sets the specified writer to
// the logger.
func SetWriter(w rogu.Writer) *rogu.Logger {
	return defaultLogger.SetWriter(w)
}

// AddWriter adds another writer to the logger.
func AddWriter(w rogu.Writer) *rogu.Logger {
	return defaultLogger.AddWriter(w)
}

// SetLevel sets the minum log leven which
// will be written.
func SetLevel(lvl level.Level) *rogu.Logger {
	return defaultLogger.SetLevel(lvl)
}

// SetCaller enabled or disables attaching the
// caller file and line to the event.
func SetCaller(enable bool) *rogu.Logger {
	return defaultLogger.SetCaller(enable)
}

// Copy creates and returns a copy of the Logger.
func Copy() *rogu.Logger {
	return defaultLogger.Copy()
}

// Tagged creates a copy of the Logger, sets the
// given tag and returns it.
func Tagged(tag string) *rogu.Logger {
	return defaultLogger.Tagged(tag)
}

// Trace creates a new log Event with level trace.
func Trace() *rogu.Event {
	return defaultLogger.Trace()
}

// Trace creates a new log Event with level debug.
func Debug() *rogu.Event {
	return defaultLogger.Debug()
}

// Trace creates a new log Event with info.
func Info() *rogu.Event {
	return defaultLogger.Info()
}

// Trace creates a new log Event with level warn.
func Warn() *rogu.Event {
	return defaultLogger.Warn()
}

// Trace creates a new log Event with level error.
func Error() *rogu.Event {
	return defaultLogger.Error()
}

// Trace creates a new log Event with level fatal.
//
// When commited, the programm will exit with exit
// code 1.
func Fatal() *rogu.Event {
	return defaultLogger.Fatal()
}

// Trace creates a new log Event with level panic.
//
// When commited, the program will panic at the
// called point.
func Panic() *rogu.Event {
	return defaultLogger.Panic()
}

package level

import "strings"

// Level specifies a log level.
type Level int8

const (
	Off Level = iota

	Panic
	Fatal
	Error
	Warn
	Info
	Debug
	Trace

	All
)

func (l Level) String() string {
	switch l {
	case Panic:
		return "panic"
	case Fatal:
		return "fatal"
	case Error:
		return "error"
	case Warn:
		return "warn"
	case Info:
		return "info"
	case Debug:
		return "debug"
	case Trace:
		return "trace"
	}

	return ""
}

// LevelFromString tries to get a Level from the
// given string.
//
// ok is false if no level could be matched
// with the passed stirng.
func LevelFromString(v string) (lvl Level, ok bool) {
	ok = true

	switch strings.ToLower(v) {
	case "panic":
		lvl = Panic
	case "fatal":
		lvl = Fatal
	case "error":
		lvl = Error
	case "warn":
		lvl = Warn
	case "info":
		lvl = Info
	case "debug":
		lvl = Debug
	case "trace":
		lvl = Trace
	default:
		ok = false
	}

	return lvl, ok
}

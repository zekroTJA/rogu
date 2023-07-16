package level

import (
	"strings"
)

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
	if v == "" {
		return Off, false
	}

	if lvl, ok = fromDigit(v); ok {
		return lvl, ok
	}

	ok = true

	switch strings.ToLower(v) {
	case "p", "pnc", "panic":
		lvl = Panic
	case "f", "ftl", "fatal":
		lvl = Fatal
	case "e", "err", "error":
		lvl = Error
	case "w", "wrn", "warn":
		lvl = Warn
	case "i", "inf", "info":
		lvl = Info
	case "d", "dbg", "debug":
		lvl = Debug
	case "t", "trc", "trace":
		lvl = Trace
	default:
		ok = false
	}

	return lvl, ok
}

func fromDigit(v string) (lvl Level, ok bool) {
	if len(v) != 1 {
		return Off, false
	}

	c := v[0]
	if c < '0' || c > '9' {
		return Off, false
	}

	return Level(c) - 48, true
}

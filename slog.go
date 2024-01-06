package rogu

import (
	"context"
	"log/slog"

	"github.com/zekrotja/rogu/level"
)

const internalErrorKey = "__internal_error_key"

var _ slog.Handler = (*logger)(nil)
var _ slog.Handler = (*Event)(nil)

func (t *logger) Enabled(_ context.Context, lvl slog.Level) bool {
	return toRoguLevel(lvl) <= t.lvl
}

func (t *logger) WithGroup(name string) slog.Handler {
	return t.Tagged(name)
}

func (t *logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return t.WithLevel(level.All).WithAttrs(attrs)
}

func (t *logger) Handle(ctx context.Context, rec slog.Record) error {
	return t.WithLevel(toRoguLevel(rec.Level)).Handle(ctx, rec)
}

func (t *taggedLogger) Handle(ctx context.Context, rec slog.Record) error {
	return t.WithLevel(toRoguLevel(rec.Level)).Handle(ctx, rec)
}

func (t *Event) Enabled(_ context.Context, lvl slog.Level) bool {
	return toRoguLevel(lvl) <= t.lvl
}

func (t *Event) WithGroup(name string) slog.Handler {
	return t.Tag(name)
}

func (t *Event) WithAttrs(attrs []slog.Attr) slog.Handler {
	for _, a := range attrs {
		if a.Key == internalErrorKey {
			t.Err(a.Value.Any().(error))
		} else {
			t.Field(a.Key, a.Value.Any())
		}
	}
	return t
}

func (t *Event) Handle(_ context.Context, rec slog.Record) error {
	t.lvl = toRoguLevel(rec.Level)

	rec.Attrs(func(a slog.Attr) bool {
		if a.Key == internalErrorKey {
			t.Err(a.Value.Any().(error))
		} else {
			t.Field(a.Key, a.Value.Any())
		}
		return true
	})

	return t.Msg(rec.Message)
}

func ErrorAttr(err error) slog.Attr {
	return slog.Attr{
		Key:   internalErrorKey,
		Value: slog.AnyValue(err),
	}
}

// ---------------------------------------------------------------------

func toRoguLevel(lvl slog.Level) level.Level {
	switch lvl {
	case slog.LevelDebug:
		return level.Debug
	case slog.LevelInfo:
		return level.Info
	case slog.LevelWarn:
		return level.Warn
	case slog.LevelError:
		return level.Error
	}

	return level.Off
}

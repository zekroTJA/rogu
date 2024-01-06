package main

import (
	"errors"
	"log/slog"

	"github.com/zekrotja/rogu"
)

func main() {
	roguLogger := slog.New(rogu.NewLogger(rogu.NewPrettyWriter()))
	slog.SetDefault(roguLogger)

	slog.Info("hello",
		"foo", "bar",
		"n", 5,
		"slice", []string{"a", "b", "c"},
		"group", slog.Group("a", "b", "c"))

	slog.With("foo", "bar").With("n", 7).Info("hello with inner args")

	dbLogger := slog.Default().WithGroup("Database")
	dbLogger.Info("Database initialized")

	wsLogger := slog.Default().WithGroup("WebServer")
	wsLogger.Error("Failed starting web server", rogu.ErrorAttr(errors.New("invalid host address")))
}

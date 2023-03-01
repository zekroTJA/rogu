package main

import (
	"errors"

	"github.com/zekrotja/rogu/log"
)

func main() {
	dbLogger := log.Tagged("Database")
	dbLogger.Info().Msg("Database initialized")

	cacheLogger := log.Tagged("Cache")
	cacheLogger.Info().Msg("Cache initialized")

	wsLogger := log.Tagged("WebServer")
	wsLogger.Error().
		Err(errors.New("invalid host address")).
		Msg("Failed starting web server")
}

package main

import (
	"errors"

	"github.com/zekrotja/rogu/level"
	"github.com/zekrotja/rogu/log"
)

func main() {
	log.SetLevel(level.All)

	log.Info().Msg("Look, this is an information!")
	log.Debug().Fields(
		"id", "ce539bd6-fd82-48a2-a7e5-d7a5eb199188",
		"counter", 78,
		"params", []interface{}{"foo", "bar", 123},
	).Msg("Some fields!")
	log.Debug().Fields(
		"a_map", map[any]any{"a": 1, 1234: "bar", "bazz": []any{5, 6, 7}},
	).Msg("Some map fields!")
	log.Error().Err(errors.New("some error")).Msg("Oh no")
	log.Trace().Caller().Msg("Here")

	log.Warn().Tag("tag 1").Msg("Uh oh")
	log.Info().Tag("bar bazz").Msg("Look, another tag!")
}

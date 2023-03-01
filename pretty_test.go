package rogu

import (
	"io"
	"testing"
)

func BenchmarkPrettyWriter(b *testing.B) {
	l := NewLogger()
	l.SetWriter(NewPrettyWriter(io.Discard))

	b.Run("single-message", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Info().Msg("bench")
		}
	})

	b.Run("fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			l.Info().Field("str", "str").Msg("bench")
		}
	})
}

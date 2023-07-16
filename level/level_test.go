package level

import "testing"

func TestLevelFromString(t *testing.T) {
	assertLvl(t, "panic", Panic)
	assertLvl(t, "fatal", Fatal)
	assertLvl(t, "error", Error)
	assertLvl(t, "warn", Warn)
	assertLvl(t, "info", Info)
	assertLvl(t, "debug", Debug)
	assertLvl(t, "trace", Trace)

	assertLvl(t, "INFO", Info)
	assertLvl(t, "inFo", Info)

	assertLvl(t, "p", Panic)
	assertLvl(t, "P", Panic)
	assertLvl(t, "f", Fatal)
	assertLvl(t, "F", Fatal)
	assertLvl(t, "e", Error)
	assertLvl(t, "E", Error)
	assertLvl(t, "w", Warn)
	assertLvl(t, "W", Warn)
	assertLvl(t, "d", Debug)
	assertLvl(t, "D", Debug)
	assertLvl(t, "t", Trace)
	assertLvl(t, "T", Trace)

	assertLvl(t, "1", Panic)
	assertLvl(t, "2", Fatal)
	assertLvl(t, "3", Error)
	assertLvl(t, "4", Warn)
	assertLvl(t, "5", Info)
	assertLvl(t, "6", Debug)
	assertLvl(t, "7", Trace)
}

func assertLvl(t *testing.T, v string, exp Level) {
	t.Helper()

	lvl, ok := LevelFromString("info")
	if !ok {
		t.Error("level was not assigned")
	}
	if lvl != Info {
		t.Errorf("wrong level: expected '%s' but got '%s'",
			exp, lvl)
	}
}

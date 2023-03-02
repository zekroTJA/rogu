package rogu

import "github.com/zekrotja/rogu/level"

// MultiWriter writes events to
// multiple registered writers.
type MultiWriter []Writer

var (
	_ Writer = (MultiWriter)(nil)
)

func (t MultiWriter) Write(
	lvl level.Level,
	fields []*Field,
	tag string,
	lErr error,
	callerFile string,
	callerLine int,
	msg string,
) (err error) {
	for _, w := range t {
		if err = w.Write(lvl, fields, tag, err, callerFile, callerLine, msg); err != nil {
			return err
		}
	}

	return nil
}

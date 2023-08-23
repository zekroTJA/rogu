package rogu

import "github.com/zekrotja/rogu/level"

// MultiWriter writes events to
// multiple registered writers.
type MultiWriter []Writer

var (
	_ Writer = (MultiWriter)(nil)
	_ Closer = (MultiWriter)(nil)
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
		if err = w.Write(lvl, fields, tag, lErr, callerFile, callerLine, msg); err != nil {
			return err
		}
	}

	return nil
}

// Close closes the set writers or all writers that
// are added to the logger and which are closable.
func (t MultiWriter) Close() error {
	for _, w := range t {
		if closer, ok := w.(Closer); ok {
			if err := closer.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

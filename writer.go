package rogu

import "github.com/zekrotja/rogu/level"

// Writer takes log entry components and
// writes them somewhere.
type Writer interface {
	Write(
		lvl level.Level,
		fields []*Field,
		tag string,
		err error,
		errFormat string,
		callerFile string,
		callerLine int,
		msg string,
	) error
}

// Closer is used to close stuff. 🤯
type Closer interface {
	Close() error
}

package rogu

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/zekrotja/rogu/level"
)

// JsonWriter implements Writer for JSON
// formatted entry output.
type JsonWriter struct {
	writeMtx sync.Mutex

	Output     io.Writer
	TimeFormat string
}

var (
	_ Writer = (*JsonWriter)(nil)
	_ Closer = (*JsonWriter)(nil)
)

// NewJsonWriter returns a new JsonWriter
// with the passed target output writers.
func NewJsonWriter(outputs ...io.Writer) *JsonWriter {
	var t JsonWriter

	if len(outputs) == 0 {
		t.Output = os.Stdout
	} else if len(outputs) == 1 {
		t.Output = outputs[0]
	} else {
		t.Output = io.MultiWriter(outputs...)
	}

	t.TimeFormat = time.RFC3339

	return &t
}

type field struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
}

type caller struct {
	File string `json:"file,omitempty"`
	Line int    `json:"line,omitempty"`
}

type entry struct {
	Timestamp string      `json:"timestamp,omitempty"`
	Level     level.Level `json:"level"`
	LevelStr  string      `json:"level_string"`
	Tag       string      `json:"tag,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Fields    []field     `json:"tags,omitempty"`
	Caller    caller      `json:"caller,omitempty"`
}

func (t *JsonWriter) Write(
	lvl level.Level,
	fields []*Field,
	tag string,
	lErr error,
	lErrFormat string,
	callerFile string,
	callerLine int,
	msg string,
) (err error) {
	var e entry

	e.Level = lvl
	e.LevelStr = lvl.String()
	e.Tag = tag
	e.Message = msg

	if t.TimeFormat != "" {
		e.Timestamp = time.Now().Format(t.TimeFormat)
	}

	if lErr != nil {
		if lErrFormat != "" {
			e.Error = fmt.Sprintf(lErrFormat, lErr)
		} else {
			e.Error = lErr.Error()
		}
	}

	if len(fields) > 0 {
		e.Fields = make([]field, 0, len(fields))
		for _, f := range fields {
			e.Fields = append(e.Fields, field{
				Key:   f.Key,
				Value: f.Val,
			})
		}
	}

	if callerFile != "" {
		e.Caller = caller{
			File: callerFile,
			Line: callerLine,
		}
	}

	t.writeMtx.Lock()
	defer t.writeMtx.Unlock()
	return json.NewEncoder(t.Output).Encode(e)
}

func (t *JsonWriter) Close() error {
	if c, ok := t.Output.(Closer); ok {
		return c.Close()
	}
	return nil
}

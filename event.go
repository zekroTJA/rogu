package rogu

import (
	"fmt"
	"reflect"

	"github.com/zekrotja/rogu/level"
)

var eventPool = newSafePool(func() *Event {
	return &Event{
		fields: make([]*Field, 0, 10),
	}
})

var fieldsPool = newSafePool(func() *Field {
	return &Field{}
})

// Field hols a hey-value pair of any type.
type Field struct {
	Key any
	Val any

	valueKind reflect.Kind
}

func (t *Field) Reset() {
	t.Key = nil
	t.Val = nil
}

// Event is used to build and send
// log messages.
type Event struct {
	lvl    level.Level
	fields []*Field
	tag    string
	err    error
	caller bool

	l *Logger
}

func (t *Event) Reset() {
	t.l = nil
	t.fields = t.fields[:0]
	t.lvl = 0
	t.tag = ""
	t.err = nil
}

func newEvent(l *Logger, lvl level.Level) *Event {
	e := eventPool.Get()
	e.l = l
	e.lvl = lvl
	return e
}

// Tag sets the tag of the event.
//
// This will overwtite the tag set by the logger.
func (t *Event) Tag(tag string) *Event {
	t.tag = tag
	return t
}

// Fields adds passed value alternating
// as keys and values to the events fields.
//
// When an odd number of values is passed,
// the last field's value will be nil.
//
// Example:
//   rogu.Info().Fields(
//       "name", "Bob",
//       "age", 24,
//       "hobbies", []stirng{"biking", "football", "gaming"},
//   )
func (t *Event) Fields(kv ...any) *Event {
	if len(kv) == 0 {
		return t
	}

	var f *Field
	for i, korv := range kv {
		if i%2 == 0 {
			// i is even, so it must be a key
			f = fieldsPool.Get()
			f.Key = korv
		} else {
			// i is odd, so it must be a value
			f.Val = korv
			t.fields = append(t.fields, f)
		}
	}

	// When the last key does not have any value,
	// commit the field with `nil` as value.
	if len(kv)%2 == 1 {
		t.fields = append(t.fields, f)
	}

	return t
}

// Field adds a single key-value field.
func (t *Event) Field(key, value any) *Event {
	f := fieldsPool.Get()
	f.Key = key
	f.Val = value
	t.fields = append(t.fields, f)
	return t
}

// Err sets an error value to the event.
func (t *Event) Err(err error) *Event {
	t.err = err
	return t
}

// Caller adds the current file and line
// to the event.
func (t *Event) Caller() *Event {
	t.caller = true
	return t
}

// Msg commits the event to the writer with
// the given message string returning an
// error when the log writing failed.
func (t *Event) Msg(v string) error {
	if t.l == nil {
		return nil
	}

	err := t.l.write(t, v)
	if err != nil {
		// Only give back the event and fields to
		// the event pools if the writing was succesfull.
		// This way, the send can safely be retried.
		t.giveBack()
	}
	return err
}

// Msgf is an alias for Msg with a format and
// given values.
func (t *Event) Msgf(format string, args ...any) error {
	return t.Msg(fmt.Sprintf(format, args...))
}

// Send commits the event without any message.
func (t *Event) Send() error {
	return t.Msg("")
}

// Discard simply throws away the event and gives
// back the used resources.
//
// This should always be called when an event
// will not be commited.
func (t *Event) Discard() {
	t.giveBack()
}

func (t *Event) giveBack() {
	for _, f := range t.fields {
		fieldsPool.Put(f)
	}
	eventPool.Put(t)
}

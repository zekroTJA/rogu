package rogu

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-colorable"
	"github.com/zekrotja/rogu/level"
)

// PrettyWriter implements Writer for human readable,
// colorful, structured console output.
//
// You can set NoColor to true to supress colorful
// formatting.
//
// With setting TimeFormat you specify the format of
// the timestamp and time.Time field values. When
// TimeFormat is set to an empty string, no
// timestamp will be printed.
//
// If you want to alter the style of the output,
// feel free to set custom definitions for
// the defined styles.
type PrettyWriter struct {
	Output io.Writer

	NoColor    bool
	TimeFormat string

	StyleTimestamp          lipgloss.Style
	StyleLevelPanic         lipgloss.Style
	StyleLevelFatal         lipgloss.Style
	StyleLevelError         lipgloss.Style
	StyleLevelWarn          lipgloss.Style
	StyleLevelInfo          lipgloss.Style
	StyleLevelDebug         lipgloss.Style
	StyleLevelTrace         lipgloss.Style
	StyleCaller             lipgloss.Style
	StyleTag                lipgloss.Style
	StyleFieldKey           lipgloss.Style
	StyleFieldValue         lipgloss.Style
	StyleFieldMultipleKey   lipgloss.Style
	StyleFieldMultipleValue lipgloss.Style
	StyleFieldErrorKey      lipgloss.Style
	StyleFieldErrorValue    lipgloss.Style
	StyleMessage            lipgloss.Style
}

var (
	_ Writer = (*PrettyWriter)(nil)
	_ Closer = (*PrettyWriter)(nil)
)

// NewPrettyWriter returns a new instance of PrettyWriter
// with the given output writers. When no writers are
// specified, os.Stdout will be used.
func NewPrettyWriter(outputs ...io.Writer) PrettyWriter {
	var t PrettyWriter

	if len(outputs) == 0 {
		t.Output = colorable.NewColorable(os.Stdout)
	} else if len(outputs) == 1 {
		t.Output = colorWriter(outputs[0])
	} else {
		for i := range outputs {
			outputs[i] = colorWriter(outputs[i])
		}
		t.Output = io.MultiWriter(outputs...)
	}

	t.TimeFormat = time.RFC3339

	t.StyleTimestamp = lipgloss.NewStyle().
		MarginRight(1).
		Foreground(lipgloss.Color("245"))

	t.StyleLevelPanic = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("201"))
	t.StyleLevelFatal = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("198"))
	t.StyleLevelError = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("196"))
	t.StyleLevelWarn = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("220"))
	t.StyleLevelInfo = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("46"))
	t.StyleLevelDebug = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("214"))
	t.StyleLevelTrace = lipgloss.NewStyle().
		MarginRight(1).
		Width(5).
		Foreground(lipgloss.Color("31"))

	t.StyleCaller = lipgloss.NewStyle().
		MarginRight(1).
		Width(18).
		Foreground(lipgloss.Color("244"))

	t.StyleTag = lipgloss.NewStyle().
		MarginRight(1).
		Width(10).
		Bold(true).
		Foreground(lipgloss.Color("45"))

	t.StyleFieldKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))
	t.StyleFieldValue = lipgloss.NewStyle().
		MarginRight(1)
	t.StyleFieldMultipleKey = t.StyleFieldKey.Copy().
		MarginTop(1).
		MarginLeft(5)
	t.StyleFieldMultipleValue = t.StyleFieldValue.Copy().
		MarginTop(1).
		MarginLeft(5).
		PaddingLeft(1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("237"))
	t.StyleFieldErrorKey = t.StyleFieldKey.Copy()
	t.StyleFieldErrorValue = lipgloss.NewStyle().
		Foreground(lipgloss.Color("160")).
		MarginRight(1)

	t.StyleMessage = lipgloss.NewStyle().
		MarginRight(1)

	return t
}

func (t PrettyWriter) Write(
	lvl level.Level,
	fields []*Field,
	tag string,
	lErr error,
	callerFile string,
	callerLine int,
	msg string,
) (err error) {
	// -- Timestamp

	if t.TimeFormat != "" {
		now := time.Now().Format(t.TimeFormat)
		if err = t.writeFormatted(now, t.StyleTimestamp); err != nil {
			return err
		}
	}

	// -- Level

	if err = t.writeLvl(lvl); err != nil {
		return err
	}

	// -- Caller

	if callerFile != "" {
		err = t.writeFormatted(t.formatCaller(callerFile, callerLine), t.StyleCaller)
		if err != nil {
			return err
		}
	}

	// -- Tag

	if tag != "" {
		if err = t.writeFormatted(tag, t.StyleTag); err != nil {
			return err
		}
	}

	// -- Message

	if err = t.writeFormatted(msg, t.StyleMessage); err != nil {
		return err
	}

	// -- Error

	if lErr != nil {
		t.writeErr(lErr)
	}

	// -- Fields

	if err = t.writeFields(fields); err != nil {
		return err
	}

	return t.writeString("\n")
}

func (t PrettyWriter) Close() error {
	if c, ok := t.Output.(Closer); ok {
		return c.Close()
	}
	return nil
}

func colorWriter(w io.Writer) io.Writer {
	if w == os.Stderr || w == os.Stdout {
		return colorable.NewColorable(w.(*os.File))
	}
	return colorable.NewNonColorable(w)
}

func (t PrettyWriter) write(p []byte) error {
	_, err := t.Output.Write(p)
	return err
}

func (t PrettyWriter) writeString(p string) error {
	_, err := fmt.Fprintf(t.Output, "%s", p)
	return err
}

func (t PrettyWriter) writeAny(v interface{}) error {
	_, err := fmt.Fprintf(t.Output, "%v", v)
	return err
}

func (t PrettyWriter) writeFormatted(v interface{}, style lipgloss.Style) (err error) {
	if t.NoColor {
		style = style.Copy().UnsetForeground()
	}
	return t.writeString(style.Render(fmt.Sprintf("%v", v)))
}

func (t PrettyWriter) writeLvl(lvl level.Level) (err error) {
	switch lvl {
	case level.Panic:
		err = t.writeFormatted("PANIC", t.StyleLevelPanic)
	case level.Fatal:
		err = t.writeFormatted("FATAL", t.StyleLevelFatal)
	case level.Error:
		err = t.writeFormatted("ERROR", t.StyleLevelError)
	case level.Warn:
		err = t.writeFormatted("WARN", t.StyleLevelWarn)
	case level.Info:
		err = t.writeFormatted("INFO", t.StyleLevelInfo)
	case level.Debug:
		err = t.writeFormatted("DEBUG", t.StyleLevelDebug)
	case level.Trace:
		err = t.writeFormatted("TRACE", t.StyleLevelTrace)
	}
	return err
}

func (t PrettyWriter) writeFields(fields []*Field) (err error) {
	for _, f := range fields {
		if f.Val != nil {
			f.valueKind = reflect.TypeOf(f.Val).Kind()
			if f.valueKind == reflect.Slice || f.valueKind == reflect.Map {
				continue
			}
		}

		err = t.writeFormatted(fmt.Sprintf("%v=", f.Key), t.StyleFieldKey)
		if err != nil {
			return err
		}

		if err = t.writeFormatted(t.valueString(f.Val), t.StyleFieldValue); err != nil {
			return err
		}
	}

	for _, f := range fields {
		if f.valueKind != reflect.Slice {
			continue
		}

		err = t.writeFormatted(fmt.Sprintf("%v=", f.Key), t.StyleFieldMultipleKey)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(f.Val)
		for i := 0; i < v.Len(); i++ {
			vi := v.Index(i).Interface()
			vStr := fmt.Sprintf("[%02d] %s", i, t.valueString(vi))
			if err = t.writeFormatted(vStr, t.StyleFieldMultipleValue); err != nil {
				return err
			}
		}
	}

	for _, f := range fields {
		if f.valueKind != reflect.Map {
			continue
		}

		err = t.writeFormatted(fmt.Sprintf("%v=", f.Key), t.StyleFieldMultipleKey)
		if err != nil {
			return err
		}

		v := reflect.ValueOf(f.Val)
		for _, k := range v.MapKeys() {
			vi := v.MapIndex(k).Interface()
			vStr := fmt.Sprintf("[%v] %s", t.valueString(k.Interface()), t.valueString(vi))
			if err = t.writeFormatted(vStr, t.StyleFieldMultipleValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t PrettyWriter) writeErr(lerr error) (err error) {
	if err = t.writeFormatted("error=", t.StyleFieldErrorKey); err != nil {
		return err
	}
	return t.writeFormatted(fmt.Sprintf("\"%s\"", lerr), t.StyleFieldErrorValue)
}

func (t PrettyWriter) valueString(v interface{}) string {
	switch vt := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", vt)
	case error:
		return fmt.Sprintf("\"%s\"", vt.Error())
	case time.Duration:
		return fmt.Sprintf("%s", vt)
	case time.Time:
		if t.TimeFormat != "" {
			return t.valueString(vt.Format(t.TimeFormat))
		}
		return t.valueString(vt)
	case interface{ String() string }:
		return fmt.Sprintf("\"%s\"", vt.String())
	}

	return fmt.Sprintf("%v", v)
}

func (t PrettyWriter) formatCaller(file string, line int) string {
	fname := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	maxFName := t.StyleCaller.GetWidth() - 2
	if len(fname) > maxFName {
		fname = "…" + fname[len(fname)-maxFName+1:]
	}

	return fmt.Sprintf("<%s>", fname)
}

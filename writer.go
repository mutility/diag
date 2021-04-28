package diag

import (
	"fmt"
	"io"
)

// NewWriter creates an Interface wrapper for an io.Writer. It will write
// Error and Warning messages to w, and discard Debug messages.
func NewWriter(w io.Writer) *wrap {
	return &wrap{io.Discard, w, w, w}
}

// NewWriterDebug creates an Interface wrapper for an io.Writer. It will write
// Error, Warning and Debug messages to w.
func NewWriterDebug(w io.Writer) *wrap {
	return &wrap{w, w, w, w}
}

// NewWriters creates an Interface wrapper for io.Writers. It will write Error,
// Warning/Print and Debug messages to their respective streams.
func NewWriters(errors, warnings, debugs io.Writer) *wrap {
	return NewWriters4(errors, warnings, warnings, debugs)
}

// NewWriters4 creates an Interface wrapper for io.Writers. It will write Error,
// Warning, Print and Debug messages to their respective streams.
func NewWriters4(errors, warnings, prints, debugs io.Writer) *wrap {
	return &wrap{wd: debugs, wp: prints, ww: warnings, we: errors}
}

type wrap struct {
	wd, wp, ww, we io.Writer
}

func (w *wrap) Debug(a ...interface{}) {
	fmt.Fprintln(w.wd, a...)
}

func (w *wrap) Print(a ...interface{}) {
	fmt.Fprintln(w.wp, a...)
}

func (w *wrap) Warning(a ...interface{}) {
	fmt.Fprintln(w.ww, a...)
}

func (w *wrap) Error(a ...interface{}) {
	fmt.Fprintln(w.we, a...)
}

// NewPrefixed returns a writer that prefixes each write with the specified
// prefix. This is useful to create differentiations for a single stream, e.g.:
//
//     log := NewWriters(NewPrefixed(w, "E:"), NewPrefixed(w, "W:"), io.Discard)
//
func NewPrefixed(w io.Writer, prefix string) *prefixWriter {
	return &prefixWriter{w, prefix}
}

type prefixWriter struct {
	w io.Writer
	p string
}

func (w *prefixWriter) Write(b []byte) (int, error) {
	var err error
	if len(b) > 0 {
		_, err = fmt.Fprintf(w.w, "%s %s", w.p, b)
	}
	return len(b), err
}

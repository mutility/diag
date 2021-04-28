// Package diag contains interfaces and utility functions to ease production
// of diagnostics. Implementations are suggested to provide at least the
// methods of diag.Interface, as it will be common for entry points to receive
// a diag.Interface. Additional methods will be leveraged where available.
//
// Typical use in a function that wants to provide diagnostics looks like this:
//
//     func Foo(log diag.Interface) {
// 	       diag.Debugf(log, "Hello %s!", "World")
//     }
//
// It's also okay to accept a diag.Debugger, diag.Errorer, or diag.Warninger
// when Foo and what it calls will use only use a subset of the capabilities.
//
// New() enables a trivial implementation around existing io.Writers, such as
// os.Stdout, os.Stderr, etc. This is useful for main or testing packages.
package diag

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type (
	Debugger  interface{ Debug(...interface{}) }
	Debugfer  interface{ Debugf(string, ...interface{}) }
	Printer   interface{ Print(...interface{}) }
	Printfer  interface{ Printf(string, ...interface{}) }
	Errorer   interface{ Error(...interface{}) }
	Errorfer  interface{ Errorf(string, ...interface{}) }
	ErrorAter interface {
		ErrorAt(string, int, int, ...interface{})
	}
	ErrorAtfer interface {
		ErrorAtf(string, int, int, string, ...interface{})
	}
	Warninger   interface{ Warning(...interface{}) }
	Warningfer  interface{ Warningf(string, ...interface{}) }
	WarningAter interface {
		WarningAt(string, int, int, ...interface{})
	}
	WarningAtfer interface {
		WarningAtf(string, int, int, string, ...interface{})
	}
	Grouper interface {
		Group(string, func(Interface))
	}
	GroupContexter interface {
		GroupContext(string, func(Context))
	}
	ValueMasker interface{ MaskValue(string) }
)

// Interface includes the core diagnostic methods. All functions in diag
// can function on top of these.
type Interface interface {
	Debugger
	Printer
	Errorer
	Warninger
}

// Context merges Interface and context.Context, enabling a single parameter to
// serve both purposes.
type Context interface {
	Interface
	context.Context
}

// FullInterface includes all diagnostic methods. This interface may be
// extended at any point, so using it confers no compatibility guarantees.
// It should be used primarily for testing if your interface is complete.
//
//     func TestFullDiagInterface(t *testing.T) {
//         var impl interface{} = (*myimpl)(nil)
//         if _, ok := impl.(diag.FullInterface); !ok {
//	           t.Error("myimpl doesn't implement diag.FullInterface")
//         }
//     }
//
// Using it any other way risks breaking your code rather than just your tests.
// For instance, the following will yield a build error for any missing methods.
//
//     var _ diag.FullInterface = (*myimpl)(nil)
//
type FullInterface interface {
	Interface
	Debugfer
	Errorfer
	ErrorAter
	ErrorAtfer
	Grouper  // added:1.2
	Printer  // added:1.1
	Printfer // added:1.1
	Warningfer
	WarningAter
	WarningAtfer
	ValueMasker
}

// WithContext creates an Interface wrapper with a Context.
func WithContext(ctx context.Context, i Interface) Context {
	return &wrapContext{ctx, i}
}

type wrapContext struct {
	context.Context
	Interface
}

// Debug outputs a debug message, unless d is nil.
func Debug(d Debugger, a ...interface{}) {
	if d != nil {
		if h := thelper(d); h != nil {
			h()
		}
		m := mask(d)
		d.Debug(m.Args(a)...)
	}
}

// Debugf outputs a formatted debug message, unless d is nil.
func Debugf(d Debugger, format string, a ...interface{}) {
	if h := thelper(d); h != nil {
		h()
	}
	if df, ok := d.(Debugfer); ok {
		m := mask(d)
		df.Debugf(m.Format(format), m.Args(a)...)
	} else if d != nil {
		m := mask(d)
		d.Debug(fmt.Sprintf(m.Format(format), m.Args(a)...))
	}
}

// Print outputs a message, unless p is nil.
//
// "Ideally" p would be a Printer instead of an Interface, but it was added late.
func Print(p Interface, a ...interface{}) {
	if p, ok := p.(Printer); ok {
		if h := thelper(p); h != nil {
			h()
		}
		p.Print(mask(p).Args(a)...)
	}
}

// Printf outputs a formatted message, unless p is nil.
//
// "Ideally" p would be a Printer instead of an Interface, but it was added late.
func Printf(p Interface, format string, a ...interface{}) {
	if h := thelper(p); h != nil {
		h()
	}
	if pf, ok := p.(Printfer); ok {
		m := mask(p)
		pf.Printf(m.Format(format), m.Args(a)...)
	} else if p, ok := p.(Printer); ok {
		m := mask(p)
		p.Print(fmt.Sprintf(m.Format(format), m.Args(a)...))
	}
}

// Error outputs an error message, unless e is nil.
func Error(e Errorer, a ...interface{}) {
	if e != nil {
		if h := thelper(e); h != nil {
			h()
		}
		e.Error(mask(e).Args(a)...)
	}
}

// Errorf outputs a formatted error message, unless e is nil.
func Errorf(e Errorer, format string, a ...interface{}) {
	if h := thelper(e); h != nil {
		h()
	}
	if ef, ok := e.(Errorfer); ok {
		m := mask(e)
		ef.Errorf(m.Format(format), m.Args(a)...)
	} else if e != nil {
		m := mask(e)
		e.Error(fmt.Sprintf(m.Format(format), m.Args(a)...))
	}
}

// ErrorAt outputs an error message with location, unless e is nil.
func ErrorAt(e Errorer, file string, line, col int, a ...interface{}) {
	if h := thelper(e); h != nil {
		h()
	}
	if ea, ok := e.(ErrorAter); ok {
		ea.ErrorAt(file, line, col, mask(e).Args(a)...)
	} else if ef, ok := e.(ErrorAtfer); ok {
		ef.ErrorAtf(file, line, col, "%s", fmt.Sprint(mask(e).Args(a)...))
	} else if e != nil {
		e.Error(fillAt(file, line, col, mask(e).Args(a))...)
	}
}

// ErrorAtf outputs a formatted error message with location, unless e is nil.
func ErrorAtf(e Errorer, file string, line, col int, format string, a ...interface{}) {
	if h := thelper(e); h != nil {
		h()
	}
	if eaf, ok := e.(ErrorAtfer); ok {
		m := mask(e)
		eaf.ErrorAtf(file, line, col, m.Format(format), m.Args(a)...)
	} else if ea, ok := e.(ErrorAter); ok {
		m := mask(e)
		ea.ErrorAt(file, line, col, fmt.Sprintf(m.Format(format), m.Args(a)...))
	} else if ef, ok := e.(Errorfer); ok {
		m := mask(e)
		ef.Errorf(fillAtf(file, line, col, m.Format(format)), m.Args(a)...)
	} else if e != nil {
		m := mask(e)
		e.Error(fmt.Sprintf(fillAtf(file, line, col, m.Format(format)), m.Args(a)...))
	}
}

// Warning outputs an warning message, unless w is nil.
func Warning(w Warninger, a ...interface{}) {
	if w != nil {
		if h := thelper(w); h != nil {
			h()
		}
		w.Warning(mask(w).Args(a)...)
	}
}

// Warningf outputs a formatted warning message, unless w is nil.
func Warningf(w Warninger, format string, a ...interface{}) {
	if h := thelper(w); h != nil {
		h()
	}
	if wf, ok := w.(Warningfer); ok {
		m := mask(w)
		wf.Warningf(m.Format(format), m.Args(a)...)
	} else if w != nil {
		m := mask(w)
		w.Warning(fmt.Sprintf(m.Format(format), m.Args(a)...))
	}
}

// WarningAt outputs an warning message with location, unless w is nil.
func WarningAt(w Warninger, file string, line, col int, a ...interface{}) {
	if h := thelper(w); h != nil {
		h()
	}
	if wa, ok := w.(WarningAter); ok {
		wa.WarningAt(file, line, col, mask(w).Args(a)...)
	} else if wf, ok := w.(WarningAtfer); ok {
		wf.WarningAtf(file, line, col, "%s", fmt.Sprint(mask(w).Args(a)...))
	} else if w != nil {
		w.Warning(fillAt(file, line, col, mask(w).Args(a))...)
	}
}

// WarningAtf outputs a formatted warning message with location, unless w is nil.
func WarningAtf(w Warninger, file string, line, col int, format string, a ...interface{}) {
	if h := thelper(w); h != nil {
		h()
	}
	if waf, ok := w.(WarningAtfer); ok {
		m := mask(w)
		waf.WarningAtf(file, line, col, m.Format(format), m.Args(a)...)
	} else if wa, ok := w.(WarningAter); ok {
		m := mask(w)
		wa.WarningAt(file, line, col, fmt.Sprintf(m.Format(format), m.Args(a)...))
	} else if wf, ok := w.(Warningfer); ok {
		m := mask(w)
		wf.Warningf(fillAtf(file, line, col, m.Format(format)), m.Args(a)...)
	} else if w != nil {
		m := mask(w)
		w.Warning(fmt.Sprintf(fillAtf(file, line, col, m.Format(format)), m.Args(a)...))
	}
}

// MaskValue requests that instances of v are obscured from output. If d
// implements ValueMasker, it fully owns the implementation. If d does not
// implement ValueMasker, then diag will obscure non-overlapping v from string
// arguments to the various output functions. (Print, Debugf, WarningAt, etc.)
//
// Diag will not obscure filenames passed to the ...At or ...Atf variants, nor
// will it attempt to obscure arguments that combine to form a requested masked
// value.
func MaskValue(d Interface, v string) {
	if m, ok := d.(ValueMasker); ok {
		m.MaskValue(v)
	} else if d != nil {
		if maskers == nil {
			maskers = make(map[interface{}]*masker)
		}
		m := maskers[d]
		if m == nil {
			m = &masker{}
			maskers[d] = m
		}
		m.masked = append(m.masked, v, "***")
		m.repl = nil
	}
}

// FormatAtBracket returns a substring of `[{{ file }}:{{ line }}.{{ col }}]`
// It terminates the inner string at the first zero value, and returns nothing
// if file is empty.
func FormatAtBracket(file string, line, col int) string {
	if file == "" {
		return ""
	}
	loc := "[" + file
	if line != 0 {
		loc += ":" + strconv.Itoa(line)
		if col != 0 {
			loc += "." + strconv.Itoa(col)
		}
	}
	loc += "]"
	return loc
}

// FormatAt globally specifies the format used for At information with
// diag.Interfaces that don't implement ...At variants. Defaults to FallbackAt.
//
// This is intended for optional overridding by package main. The first empty
// or zero value of file, line, and col indicate the rest should also be
// ignored, but this is not enforced by diag.
//
// If you need different behaviors for warning and error, you should implement
// the ...At variants directly.
var FormatAt = FormatAtBracket

func fillAt(file string, line, col int, a []interface{}) []interface{} {
	if loc := FormatAt(file, line, col); loc != "" {
		return append([]interface{}{loc}, a...)
	}
	return a
}

func fillAtf(file string, line, col int, format string) string {
	loc := FormatAt(file, line, col)
	if loc == "" {
		return format
	}
	return strings.ReplaceAll(loc, "%", "%%") + " " + format
}

// thelper retrieves a t.Helper() method if i implements it. This allows
// diag to use t.Helper() to disappear from the logging locations.
func thelper(i interface{}) func() {
	if h, ok := i.(interface {
		Helper()
	}); ok {
		return h.Helper
	}
	return nil
}

type masker struct {
	masked []string
	repl   *strings.Replacer
}

var maskers map[interface{}]*masker

func mask(d interface{}) *masker {
	m := maskers[d]
	if m == nil || len(m.masked) == 0 {
		return nil
	}
	if m.repl == nil {
		m.repl = strings.NewReplacer(m.masked...)
	}
	return m
}

func (m *masker) Args(a []interface{}) []interface{} {
	if m == nil {
		return a
	}
	repl := m.repl
	a = append([]interface{}(nil), a...)
	for i := range a {
		if s, ok := a[i].(string); ok {
			a[i] = repl.Replace(s)
		}
	}
	return a
}

func (m *masker) Format(format string) string {
	if m == nil {
		return format
	}
	return m.repl.Replace(format)
}

package diag_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mutility/diag"
)

// TestNil verifies that diag.* functions accept nil implementations without panic.
func TestNil(t *testing.T) {
	format := "%s %v %d%v"
	args := []interface{}{"a", "b", 2, 3}
	file := "somefile.md"
	line, col := 6, 0
	for name, fn := range map[string]func(){
		"Debug":      func() { diag.Debug(nil, args...) },
		"Debugf":     func() { diag.Debugf(nil, format, args...) },
		"Print":      func() { diag.Print(nil, args...) },
		"Printf":     func() { diag.Printf(nil, format, args...) },
		"Warning":    func() { diag.Warning(nil, args...) },
		"WarningAt":  func() { diag.WarningAt(nil, file, line, col, args...) },
		"Warningf":   func() { diag.Warningf(nil, format, args...) },
		"WarningAtf": func() { diag.WarningAtf(nil, file, line, col, format, args...) },
		"Error":      func() { diag.Error(nil, args...) },
		"ErrorAt":    func() { diag.ErrorAt(nil, file, line, col, args...) },
		"Errorf":     func() { diag.Errorf(nil, format, args...) },
		"ErrorAtf":   func() { diag.ErrorAtf(nil, file, line, col, format, args...) },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Error("recovered from panic")
				}
			}()
			fn()
		})
	}
}

// TestFill ensures that ...f, ...At, and ...Atf methods are wired to the underlying
func TestFill(t *testing.T) {
	d := &fill{}
	format := "%s %v %d%v"
	args := []interface{}{"a", "b", 2, 3}
	file := "somefile.md"
	line, col := 6, 0
	for name, fn := range map[string]func() string{
		"Debug":      func() string { diag.Debug(d, args...); return d.debug() },
		"Debugf":     func() string { diag.Debugf(d, format, args...); return d.debug() },
		"Print":      func() string { diag.Print(d, args...); return d.print() },
		"Printf":     func() string { diag.Printf(d, format, args...); return d.print() },
		"Warning":    func() string { diag.Warning(d, args...); return d.warning() },
		"WarningAt":  func() string { diag.WarningAt(d, file, line, col, args...); return d.warning() },
		"Warningf":   func() string { diag.Warningf(d, format, args...); return d.warning() },
		"WarningAtf": func() string { diag.WarningAtf(d, file, line, col, format, args...); return d.warning() },
		"Error":      func() string { diag.Error(d, args...); return d.error() },
		"ErrorAt":    func() string { diag.ErrorAt(d, file, line, col, args...); return d.error() },
		"Errorf":     func() string { diag.Errorf(d, format, args...); return d.error() },
		"ErrorAtf":   func() string { diag.ErrorAtf(d, file, line, col, format, args...); return d.error() },
	} {
		t.Run(name, func(t *testing.T) {
			suffix := strings.TrimPrefix(name, "Wa")[5:]
			got := fn()
			want := map[string]string{
				"":    "a b 2 3\n",
				"f":   "a b 23\n",
				"At":  "[somefile.md:6] a b 2 3\n",
				"Atf": "[somefile.md:6] a b 23\n",
			}[suffix]
			if got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		})
	}
}

// TestFillMask ensures that ...f, ...At, and ...Atf methods are wired to the underlying
func TestFillMask(t *testing.T) {
	d := &fill{}
	diag.MaskValue(d, "abc")
	format := "%s%s%s %q"
	args := []interface{}{"a", "b", "c", "abc"}
	file := "somefile.abc"
	line, col := 6, 0
	for name, fn := range map[string]func() string{
		"Debug":      func() string { diag.Debug(d, args...); return d.debug() },
		"Debugf":     func() string { diag.Debugf(d, format, args...); return d.debug() },
		"Print":      func() string { diag.Print(d, args...); return d.print() },
		"Printf":     func() string { diag.Printf(d, format, args...); return d.print() },
		"Warning":    func() string { diag.Warning(d, args...); return d.warning() },
		"WarningAt":  func() string { diag.WarningAt(d, file, line, col, args...); return d.warning() },
		"Warningf":   func() string { diag.Warningf(d, format, args...); return d.warning() },
		"WarningAtf": func() string { diag.WarningAtf(d, file, line, col, format, args...); return d.warning() },
		"Error":      func() string { diag.Error(d, args...); return d.error() },
		"ErrorAt":    func() string { diag.ErrorAt(d, file, line, col, args...); return d.error() },
		"Errorf":     func() string { diag.Errorf(d, format, args...); return d.error() },
		"ErrorAtf":   func() string { diag.ErrorAtf(d, file, line, col, format, args...); return d.error() },
	} {
		t.Run(name, func(t *testing.T) {
			suffix := strings.TrimPrefix(name, "Wa")[5:]
			got := fn()
			want := map[string]string{
				"":    "a b c ***\n",
				"f":   "abc \"***\"\n",
				"At":  "[somefile.abc:6] a b c ***\n",
				"Atf": "[somefile.abc:6] abc \"***\"\n",
			}[suffix]
			if got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		})
	}
}

type fill struct {
	d, p, w, e string
}

func (f *fill) Debug(a ...interface{})   { f.d = fmt.Sprintln(a...) }
func (f *fill) Print(a ...interface{})   { f.p = fmt.Sprintln(a...) }
func (f *fill) Warning(a ...interface{}) { f.w = fmt.Sprintln(a...) }
func (f *fill) Error(a ...interface{})   { f.e = fmt.Sprintln(a...) }
func (f *fill) debug() string            { s := f.d; f.d = ""; return s }
func (f *fill) print() string            { s := f.p; f.p = ""; return s }
func (f *fill) warning() string          { s := f.w; f.w = ""; return s }
func (f *fill) error() string            { s := f.e; f.e = ""; return s }

// TestAt verifies the At prefixes stop at the first unspecified item
func TestAt(t *testing.T) {
	d := &fill{}
	c := &customat{}
	for _, tt := range []struct {
		why       string
		file      string
		line, col int
		args      interface{}
		want      string
		cwant     string
	}{
		{"noinfo", "", 0, 0, "args", "args\n", "[|0|0]args\n"},
		{"nofile", "", 10, 3, "args", "args\n", "[|10|3]args\n"},
		{"file", "fn.go", 0, 0, "args", "[fn.go] args\n", "[fn.go|0|0]args\n"},
		{"noline", "fn.go", 0, 3, "args", "[fn.go] args\n", "[fn.go|0|3]args\n"},
		{"line", "fn.go", 10, 0, "args", "[fn.go:10] args\n", "[fn.go|10|0]args\n"},
		{"all", "fn.go", 10, 3, "args", "[fn.go:10.3] args\n", "[fn.go|10|3]args\n"},
	} {
		t.Run(tt.why, func(t *testing.T) {
			diag.WarningAt(d, tt.file, tt.line, tt.col, tt.args)
			got := d.warning()
			if got != tt.want {
				t.Errorf("fill: got %q; want %q", got, tt.want)
			}

			diag.WarningAt(c, tt.file, tt.line, tt.col, tt.args)
			got = c.warning()
			if got != tt.cwant {
				t.Errorf("custom: got %q; want %q", got, tt.cwant)
			}
		})
	}
}

type customat struct {
	fill
}

func (c *customat) WarningAt(file string, line, col int, args ...interface{}) {
	c.w = fmt.Sprintf("[%s|%d|%d]", file, line, col) + fmt.Sprintln(args...)
}

func (c *customat) ErrorAt(file string, line, col int, args ...interface{}) {
	c.w = fmt.Sprintf("[%s|%d|%d]", file, line, col) + fmt.Sprintln(args...)
}

// TestFallback verifies which fallback is used
func TestFallback(t *testing.T) {
	var got string
	f := &hasf{&got}
	at := &hasat{&got}
	atf := &hasatf{&got}
	for d, wants := range map[diag.Interface]struct{ base, f, at, atf string }{
		f:   {"", "f", "", "f"},
		at:  {"", "", "At", "At"},
		atf: {"", "", "Atf", "Atf"}, // prefer Atf over base for At()
	} {
		t.Run(fmt.Sprint(d), func(t *testing.T) {
			test := func(name string, fn func(), want string) {
				t.Run(name, func(t *testing.T) {
					fn()
					if got != want {
						t.Errorf("called %s, want %s", got, want)
					}
				})
			}
			test("Debug", func() { diag.Debug(d, "d") }, "Debug"+wants.base)
			test("Debugf", func() { diag.Debugf(d, "d") }, "Debug"+wants.f)
			test("Print", func() { diag.Print(d, "d") }, "Print"+wants.base)
			test("Printf", func() { diag.Printf(d, "d") }, "Print"+wants.f)
			test("Warning", func() { diag.Warning(d, "d") }, "Warning"+wants.base)
			test("Warningf", func() { diag.Warningf(d, "d") }, "Warning"+wants.f)
			test("WarningAt", func() { diag.WarningAt(d, "f", 1, 2, "d") }, "Warning"+wants.at)
			test("WarningAtf", func() { diag.WarningAtf(d, "f", 1, 2, "d") }, "Warning"+wants.atf)
			test("Error", func() { diag.Error(d, "d") }, "Error"+wants.base)
			test("Errorf", func() { diag.Errorf(d, "d") }, "Error"+wants.f)
			test("ErrorAt", func() { diag.ErrorAt(d, "f", 1, 2, "d") }, "Error"+wants.at)
			test("ErrorAtf", func() { diag.ErrorAtf(d, "f", 1, 2, "d") }, "Error"+wants.atf)
		})
	}
}

type hasf struct{ called *string }

func (h *hasf) Debug(...interface{})            { *h.called = "Debug" }
func (h *hasf) Debugf(string, ...interface{})   { *h.called = "Debugf" }
func (h *hasf) Print(...interface{})            { *h.called = "Print" }
func (h *hasf) Printf(string, ...interface{})   { *h.called = "Printf" }
func (h *hasf) Warning(...interface{})          { *h.called = "Warning" }
func (h *hasf) Warningf(string, ...interface{}) { *h.called = "Warningf" }
func (h *hasf) Error(...interface{})            { *h.called = "Error" }
func (h *hasf) Errorf(string, ...interface{})   { *h.called = "Errorf" }
func (h *hasf) String() string                  { return "hasf" }

type hasat struct{ called *string }

func (h *hasat) Debug(...interface{})                       { *h.called = "Debug" }
func (h *hasat) DebugAt(string, int, int, ...interface{})   { *h.called = "DebugAt" }
func (h *hasat) Print(...interface{})                       { *h.called = "Print" }
func (h *hasat) PrintAt(string, int, int, ...interface{})   { *h.called = "PrintAt" }
func (h *hasat) Warning(...interface{})                     { *h.called = "Warning" }
func (h *hasat) WarningAt(string, int, int, ...interface{}) { *h.called = "WarningAt" }
func (h *hasat) Error(...interface{})                       { *h.called = "Error" }
func (h *hasat) ErrorAt(string, int, int, ...interface{})   { *h.called = "ErrorAt" }
func (h *hasat) String() string                             { return "hasAt" }

type hasatf struct{ called *string }

func (h *hasatf) Debug(...interface{})                                { *h.called = "Debug" }
func (h *hasatf) DebugAtf(string, int, int, string, ...interface{})   { *h.called = "DebugAtf" }
func (h *hasatf) Print(...interface{})                                { *h.called = "Print" }
func (h *hasatf) PrintAtf(string, int, int, string, ...interface{})   { *h.called = "PrintAtf" }
func (h *hasatf) Warning(...interface{})                              { *h.called = "Warning" }
func (h *hasatf) WarningAtf(string, int, int, string, ...interface{}) { *h.called = "WarningAtf" }
func (h *hasatf) Error(...interface{})                                { *h.called = "Error" }
func (h *hasatf) ErrorAtf(string, int, int, string, ...interface{})   { *h.called = "ErrorAtf" }
func (h *hasatf) String() string                                      { return "hasAtf" }

func TestPrefix(t *testing.T) {
	for want, input := range map[string]string{
		"":          "",
		"p:: abc\n": "abc\n",
	} {
		t.Run(input, func(t *testing.T) {
			sb := &strings.Builder{}
			p := diag.NewPrefixed(sb, "p::")
			n, err := p.Write([]byte(input))
			if err != nil {
				t.Error("unexpected write error:", err)
			}
			if n != len(input) {
				t.Errorf("wrote %d bytes; want %d", n, len(input))
			}
			if got := sb.String(); got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		})
	}
}

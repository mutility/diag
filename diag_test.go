package diag_test

import (
	"strings"
	"testing"

	"github.com/mutility/diag"
)

func TestWrapper(t *testing.T) {
	test := func(
		fn interface{},
		format string,
		args []interface{},
		want string,
		at ...interface{},
	) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			sb := &strings.Builder{}
			log := diag.New(sb)
			log.DebugPrefix = ""
			log.WarningPrefix = ""
			log.ErrorPrefix = ""

			switch fn := fn.(type) {

			// base
			case func(diag.Debugger, ...interface{}):
				fn(log, args...)
			case func(diag.Warninger, ...interface{}):
				fn(log, args...)
			case func(diag.Errorer, ...interface{}):
				fn(log, args...)

			// f
			case func(diag.Debugger, string, ...interface{}):
				fn(log, format, args...)
			case func(diag.Warninger, string, ...interface{}):
				fn(log, format, args...)
			case func(diag.Errorer, string, ...interface{}):
				fn(log, format, args...)

			// At
			case func(diag.Debugger, string, int, int, ...interface{}):
				fn(log, at[0].(string), at[1].(int), at[2].(int), args...)
			case func(diag.Warninger, string, int, int, ...interface{}):
				fn(log, at[0].(string), at[1].(int), at[2].(int), args...)
			case func(diag.Errorer, string, int, int, ...interface{}):
				fn(log, at[0].(string), at[1].(int), at[2].(int), args...)

			// Atf
			case func(diag.Debugger, string, int, int, string, ...interface{}):
				fn(log, at[0].(string), at[1].(int), at[2].(int), format, args...)
			case func(diag.Warninger, string, int, int, string, ...interface{}):
				fn(log, at[0].(string), at[1].(int), at[2].(int), format, args...)
			case func(diag.Errorer, string, int, int, string, ...interface{}):
				fn(log, at[0].(string), at[1].(int), at[2].(int), format, args...)

			default:
				t.Errorf("Unexpected fn type %T", fn)
			}
			if got := sb.String(); got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	}

	for want, args := range map[string][]interface{}{
		"args 0 one\n": {"args", 0, "one"},
		"args 2 3\n":   {"args", 2, 3},
	} {
		t.Run(want, func(t *testing.T) {
			t.Run("Debug", test(diag.Debug, "", args, ""))
			t.Run("Warning", test(diag.Warning, "", args, want))
			t.Run("Error", test(diag.Error, "", args, want))
		})
	}

	t.Run("f", func(t *testing.T) {
		for want, args := range map[string][]interface{}{
			"args 0 one\n": {"args %d one", 0},
			"args 2 3\n":   {"args %d %v", 2, 3},
		} {
			t.Run(want, func(t *testing.T) {
				format := args[0].(string)
				args = args[1:]
				t.Run("Debug", test(diag.Debugf, format, args, ""))
				t.Run("Warning", test(diag.Warningf, format, args, want))
				t.Run("Error", test(diag.Errorf, format, args, want))
			})
		}
	})

	t.Run("At", func(t *testing.T) {
		for want, args := range map[string][]interface{}{
			"file 0\n":          {"", 0, 0, "file", 0},
			"[fn] file 1\n":     {"fn", 0, 0, "file", 1},
			"[fn:4] file 2\n":   {"fn", 4, 0, "file", 2},
			"[fn:4.8] file 3\n": {"fn", 4, 8, "file", 3},
		} {
			t.Run(want, func(t *testing.T) {
				at := args[:3]
				args = args[3:]
				t.Run("Warning", test(diag.WarningAt, "", args, want, at...))
				t.Run("Error", test(diag.ErrorAt, "", args, want, at...))
			})
		}
	})

	t.Run("Atf", func(t *testing.T) {
		for want, args := range map[string][]interface{}{
			"file 0\n":          {"", 0, 0, "file %v", 0},
			"[fn] file 1\n":     {"fn", 0, 0, "file %d", 1},
			"[fn:4] file 2\n":   {"fn", 4, 0, "file %d", 2},
			"[fn:4.8] file 3\n": {"fn", 4, 8, "file %v", "3"},
		} {
			t.Run(want, func(t *testing.T) {
				at := args[:3]
				format := args[3].(string)
				args = args[4:]
				t.Run("Warning", test(diag.WarningAtf, format, args, want, at...))
				t.Run("Error", test(diag.ErrorAtf, format, args, want, at...))
			})
		}
	})
}

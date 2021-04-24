// package testdiag adapts a testing.TB to a diag.Interface or diag.Context.
package testdiag

import (
	"context"

	"github.com/mutility/diag"
)

type testDiag struct {
	t testlog
}

// testlog is the subset of testing.TB we need
type testlog interface {
	Helper()
	Log(...interface{})
}

// Interface returns a diag.Interface that logs to t
func Interface(tb testlog) diag.Interface {
	return testDiag{tb}
}

// Context returns a diag.Context that logs to t and uses context.Background
func Context(tb testlog) diag.Context {
	return WithContext(context.Background(), tb)
}

// Context returns a diag.Context that logs to t and uses the specified context
func WithContext(ctx context.Context, tb testlog) diag.Context {
	return diag.WithContext(ctx, Interface(tb))
}

func (d testDiag) Debug(args ...interface{})   { d.t.Helper(); d.t.Log(args...) }
func (d testDiag) Warning(args ...interface{}) { d.t.Helper(); d.t.Log(args...) }
func (d testDiag) Error(args ...interface{})   { d.t.Helper(); d.t.Log(args...) }

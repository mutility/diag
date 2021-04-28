package diag

import "context"

// Group begins a grouped section of output. If d implements Grouper, it
// owns the implementation and its behavior. If not, diag will indent lines
// output during the call to fn.
//
// It is not well-defined what happens if methods on d are called during fn.
func Group(d Interface, title string, fn func(Interface)) {
	if h := thelper(d); h != nil {
		h()
	}
	if g, ok := d.(Grouper); ok {
		g.Group(title, fn)
	} else {
		Printf(d, "%s:", title)
		fn(&grouped{d})
	}
}

// GroupContext begins a grouped section of output. If d implements
// GroupContexter, it // owns the implementation and its behavior. If not, diag
// will indent lines output during the call to fn.
//
// It is not well-defined what happens if methods on d are called during fn.
func GroupContext(d Context, title string, fn func(Context)) {
	if h := thelper(d); h != nil {
		h()
	}
	if g, ok := d.(GroupContexter); ok {
		g.GroupContext(title, fn)
	} else {
		Printf(d, "%s:", title)
		fn(&groupedctx{grouped{d}, d})
	}
}

type groupedctx struct {
	grouped
	context.Context
}

type grouped struct {
	d Interface
}

func (g *grouped) Debug(a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Debug(g.d, append([]interface{}{" "}, a...)...)
}

func (g *grouped) Debugf(format string, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Debugf(g.d, "  "+format, a...)
}

func (g *grouped) Print(a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Print(g.d, append([]interface{}{" "}, a...)...)
}

func (g *grouped) Printf(format string, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Printf(g.d, "  "+format, a...)
}

func (g *grouped) Warning(a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Warning(g.d, append([]interface{}{" "}, a...)...)
}

func (g *grouped) Warningf(format string, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Warningf(g.d, "  "+format, a...)
}

func (g *grouped) WarningAt(file string, line, col int, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	WarningAt(g.d, file, line, col, append([]interface{}{" "}, a...)...)
}

func (g *grouped) WarningAtf(file string, line, col int, format string, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	WarningAtf(g.d, file, line, col, "  "+format, a...)
}

func (g *grouped) Error(a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Error(g.d, append([]interface{}{" "}, a...)...)
}

func (g *grouped) Errorf(format string, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	Errorf(g.d, "  "+format, a...)
}

func (g *grouped) ErrorAt(file string, line, col int, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	ErrorAt(g.d, file, line, col, append([]interface{}{" "}, a...)...)
}

func (g *grouped) ErrorAtf(file string, line, col int, format string, a ...interface{}) {
	if h := thelper(g.d); h != nil {
		h()
	}
	ErrorAtf(g.d, file, line, col, "  "+format, a...)
}

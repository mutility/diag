# diag

Diag provides an interface adapter for command line applications to issue diagnostic output. The application can issue debug, warning, and error messages using wrappers similar to `fmt.Printf`.

[![CI](https://github.com/mutility/diag/actions/workflows/build.yaml/badge.svg)](https://github.com/mutility/diag/actions/workflows/build.yaml)

## Consuming diag.Interface

Import diag, and accept a `diag.Interface` or `diag.Context` in your library functions. Functions that don't accept a `context.Context` should likely accept a `diag.Interface`. Functions that otherwise would accept a `context.Context` can accept a `diag.Context` instead, or can keep them separate by accepting both a `context.Context` and a `diag.Interface`.

Then to issue debug messages, warnings, or errors, use functions from `diag`: `diag.Debug`, `diag.Warning`, `diag.Error`. These also come in a mix of `...f`, `...At`, and `...Atf` variants.

Main or test routines can use helpers from `appdiag`, `ghadiag` or `testdiag` to provide a `diag.Interface` or `diag.Context`. This trivial example shows both inline, although typically they should be split across your main and library packages.

```go
package main

import (
  "github.com/mutility/diag"
)

func main() {
    Do(diag.NewOS())
    DoContext(diag.NewOSContext(context.Background()))
}

func Do(log diag.Interface) {
  diag.Debug(log, "enter Do")
  diag.Error(log, "Do is not implemented yet")
  diag.Debug(log, "exit Do")
}

func DoContext(ctx diag.Context) {
  if err := ctx.Err(); err != nil {
    diag.Warn(ctx, "context error:", err)
  }
}
```

## Testing with diag.Interface

The `testdiag` package provides functions `Interface`, `Context`, and `WithContext` that adapt a `testing.TB` to `diag.Interface`, `diag.Context` (using `context.Background`), and `diag.Context` (using a supplied context) respectively.

If you prefer to capture and process the output, you can instead wrap a `strings.Builder` or other `io.Writer` with `diag.NewWriter` or `diag.NewWriters`. If you want prefixes, wrap the writer first with `diag.NewPrefixed`.

Alternately, the functions in `diag` politely do nothing if a nil is passed as the `diag.Interface`. (Just make sure to pass the untyped nil, not a typed nil, unless that type's implementation works with an underlying nil pointer.)

## Implementing diag.Interface

You can implement anything between `diag.Interface` and `diag.FullInterface`, and the functions in `diag` will make up the difference. As an example, you can see the `testdiag` implementation inclues only `Debug`, `Warning`, and `Error` methods that each call `tb.Log`.

## Inspiration

Diag was created to ease creation of utilities that work well both on the command line and as an action in a workflow on GitHub Actions. Most of the methods reflect what GitHub Actions make available, with trivial fallback approaches for more regular scenarios.

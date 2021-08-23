# errors

[![codecov](https://codecov.io/gh/./branch/master/graph/badge.svg)](https://codecov.io/gh/.)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/.)
[![Go Report Card](https://goreportcard.com/badge/.)](https://goreportcard.com/report/github.com/adamhassel/errors)

Package errors implements a stdlib-compatible way of wrapping more than
one error into an error chain, while supporting `errors.Is` and `errors.As` (and
obviously `Error()` and `Unwrap()`), thus being a drop in replacement for other
error types. Errors also provides a `New()` function which works like the stdlib version.

The Standard Library allows you to wrap exactly ONE error with eg
`fmt.Errorf`, and will bail out if there are more than one `%w` receiver.

`errors.Wrap` will solve that by allowing arbitrary errors to be wrapped without losing information.

A construed (well, kinda) could be a situation where a bunch of different functions are called, and if failing, returns a common error that is handled further up the stack:

```golang
if err := someFunc(); err != nil {
	// MyError is handled up the stack
	return MyError
}
```

This discards information about the actual error, though. Now, normally, you'd do something like:

```golang
return fmt.Errorf("%s: %w", err, MyError)
```

And while this wraps `MyError`, it doesn't wrap the actual error from `someFunc, which we might be interested in as well, maybe because some specific error requires specific handling:

So, if the error above was e.g. `mysql.ErrNotFound`, `errors.Is(err, mysql.ErrNotFound)` would be false, even if `errors.Is(err, MyError)` is true.

Instead, using `errors.Wrap`:

```golang
if err := someFunc(); err != nil {
	// MyError is handled explicitly up the stack
	return errors.Wrap(err, MyError)
}
```

This will make `errors.Is` return true for both, as they are now both properly wrapped.

`Unwrap` and `Wrap` both run in O(n) time, where n is the number of errors added to the chain.

## Functions

### func [New](/errors.go#L88)

`func New(text string) error`

New is a wrapper around the stdlib errors.New function

### func [Wrap](/errors.go#L73)

`func Wrap(errs ...error) error`

Wrap will wrap one or more errors into a single error chain, compatible with
`errors.As` and `errors.Is`. Note that if you're using this as a sort of `append`
analogue (`err = Wrap(err, ErrAnother)` or similar, where the result
overwrites the argument), then you should protect accordingly with appropriate
synchronization measures (e.g. a mutex), just as you would with `append`.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)

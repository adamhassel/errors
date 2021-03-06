# errors

[![codecov](https://codecov.io/gh/adamhassel/errors/branch/main/graph/badge.svg?token=Eo4O4BXZ2N)](https://codecov.io/gh/adamhassel/errors)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/.)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamhassel/errors)](https://goreportcard.com/report/github.com/adamhassel/errors)

Package errors implements a stdlib-compatible way of wrapping more than
one error into an error chain, while supporting `errors.Is` and `errors.As` (and
obviously `Error()` and `Unwrap()`), thus being a drop in replacement for other
error types. Errors also provides all the stdlib functions, so you don't need to include both.

The Standard Library allows you to wrap exactly ONE error with eg
`fmt.Errorf`, and will bail out if there are more than one `%w` receiver.

`errors.Wrap` will solve that by allowing arbitrary errors to be wrapped without losing information.

This is useful e.g. when collecting errors from several running go routines that might return a bunch of different errors.

Another slightly construed use case could be a situation where a bunch of different functions are called, and if failing, returns a common error that is handled further up the stack:

```golang
if err := someFunc(); err != nil {
	// MyError is handled up the stack
	return MyError
}
```

Up the stack, you don't want to check for any possible error. You just know that if you get `MyError`, you're handling things one way, and if you get any other error, you handle it differently. But this approach, while convenient, discards information about the actual error, though. Now, normally, you'd log it at the source at the risk of it being logged again later, or do something like:

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

### func [Wrap](/errors.go#L73)

`func Wrap(errs ...error) error`

Wrap will wrap one or more errors into a single error chain, compatible with
`errors.As` and `errors.Is`. Note that if you're using this as a sort of `append`
analogue (`err = Wrap(err, ErrAnother)` or similar, where the result
overwrites the argument), then you should protect accordingly with appropriate
synchronization measures (e.g. a mutex), just as you would with `append`.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)

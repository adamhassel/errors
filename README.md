# errors

[![Build Status](https://travis-ci.org/adamhassel/errors.svg?branch=master)](https://travis-ci.org/adamhassel/errors)
[![codecov](https://codecov.io/gh/adamhassel/errors/branch/master/graph/badge.svg)](https://codecov.io/gh/adamhassel/errors)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/adamhassel/errors)

Package errors implements a stdlib-compatible way of wrapping more than
one error into an error chain, while supporting `errors.Is` and `errors.As` (and
obviously `Error()` and `Unwrap()`), thus being a drop in replacement for other
error types. Errors also provides a `New()` function which works like the stdlib version.

The Standard Library allows you to wrap exactly ONE error with eg
`fmt.Errorf`, and will bail out if there are more than one `%!(NOVERB)w` receiver.

`errors.Wrap` will solve that by allowing arbitrary errors to be wrapped without losing information.

A construed (well, kinda) could be a situation where a bunch of different functions are called, and if failing, returns a common error that is handled further up the stack:

```go
if err := someFunc(); err != nil {
// MyError is handled up the stack
return MyError
}
```

This discards information about the actual error, though. Now, normally, you'd do something like:

```go
return fmt.Errorf("%!(NOVERB)s: %!(NOVERB)w", err, MyError)
```

And while this wraps `MyError`, it doesn't wrap the actual error from `someFunc, which we might be interested in as well, maybe because some specific error requires specific handling:

So, if the error above was e.g. `mysql.ErrNotFound`, `errors.Is(err, mysql.ErrNotFound)` would be false, even if `errors.Is(err, MyError)` is true.

Instead, using `errors.Wrap`:

```go
if err := someFunc(); err != nil {

// MyError is handled explicitly up the stack
return errors.Wrap(err, MyError)

}
```

This will make `errors.Is` return true for both, as they are now both properly wrapped.

`Unwrap` and `Wrap` both run in O(n) time, where n is the number of errors added to the chain.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)

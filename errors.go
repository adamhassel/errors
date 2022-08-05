package errors

import (
	"errors"
)

// echain is a linked list (kinda) containing errors. But since echain is an
// error itself, there's nothing stopping anyone from stuffing an echain into the
// err value, and then we're actually working with a weird mutant tree instead.
// 🤷 Either way, it's depth first, and the runtime of Unwrap and thus Is and As
// is always O(n), no matter the structure...
type echain struct {
	err  error
	next *echain
}

// Slice converts an error chain to a slice of errors
func (ec *echain) Slice() []error {
	out := make([]error, 0, 1)
	out = append(out, ec.err)
	if ec.next != nil {
		out = append(out, ec.next.Slice()...)
	}
	return out
}

// Len computes the length of an error chain
func (ec *echain) Len() int {
	var out = 1 // count the err in this link, even if it's nil
	if ec.next != nil {
		out++
		out += ec.next.Len()
	}
	return out
}

// add an error to the chain. Locked from caller.
func (ec *echain) add(err error) {
	if ec.next == nil {
		ec.next = &echain{err: err}
		return
	}
	ec.next.add(err)
}

// Error implements the error interface
func (ec *echain) Error() string {
	return ErrorFormatter(ec.Slice())
}

// As is the implementation of errors.As, ensuring the chain works
func (ec *echain) As(target interface{}) bool {
	return errors.As(ec.err, &target)
}

// Is is the implementation of errors.Is, ensuring that the chain works
func (ec *echain) Is(target error) bool {
	return errors.Is(ec.err, target)
}

// Unwrap makes sure echain can act as a golang error chain
func (ec *echain) Unwrap() (err error) {
	if ec == nil {
		return nil
	}
	if ec.next == nil {
		return ec.err
	}
	return ec.next
}

// Wrap will wrap one or more errors into a single error chain, compatible with
// errors.As, errors.Is. Note that if you're using this as a sort of `append`
// analogue (`err = Wrap(err, ErrAnother)` or similar, where the result
// overwrites the argument), then you should protect accordingly with appropriate
// synchronization measures (e.g. a mutex), just as you would with `append`.
func Wrap(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	var e error
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		if e == nil {
			e = &echain{}
		}
		e.(*echain).add(errs[i])
	}
	return e
}

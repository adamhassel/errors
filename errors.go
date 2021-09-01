package errors

import (
	"errors"
	"strings"
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

// Separator between errors in the chain when they're output together with
// `.Error()`. Can be overridden globally by setting this value, default is ": ".
var Separator = ": "

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
	if ec.next == nil {
		return ec.err.Error()
	}
	errs := make([]string, 0)

	for ec != nil {
		if ec.err != nil {
			errs = append(errs, ec.err.Error())
		}
		ec = ec.next
	}
	return strings.Join(errs, Separator)
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
func (ec *echain) Unwrap() error {
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
	if len(errs) == 1 {
		return errs[0]
	}
	e := &echain{}
	for _, err := range errs {
		e.add(err)
	}
	return e
}

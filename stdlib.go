package errors

import "errors"

// New is a wrapper around the stdlib errors.New function, so you don't need to import both
func New(text string) error {
	return errors.New(text)
}

// Is is a wrapper around the stdlib function, so you don't need to import both packages
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is a wrapper around the stdlib function, so you don't need to import both packages
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap is a wrapper around the stdlib function, so you don't need to import both packages
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

package errors

import "strings"

type ErrFmtFunc func([]error) string

// Separator between errors in the chain when they're output together with
// `.Error()`. Can be overridden globally by setting this value, default is ": ". Used by the default ErrorFormatter as a separator
var Separator = ": "

// errFmt is the default formatter
var errFmt = func(es []error) string {
	out := make([]string, 0, len(es))
	for _, e := range es {
		if e == nil {
			continue
		}
		out = append(out, e.Error())
	}
	return strings.Join(out, Separator)
}

var ErrorFormatter ErrFmtFunc = errFmt

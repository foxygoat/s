// Package errs provides error wrapping functions that allow an error to
// wrap multiple errors. The error wrapping in the Go standard library
// errors package allows an error to wrap only one error with one %w
// format verb in the format string passed to fmt.Errorf and the
// errors.Unwrap() function that returns a single error.
//
// The error type returned by the functions in this package wraps
// multiple errors so that errors.Is() and errors.As() can be used to
// query if that error is any one of the wrapped errors. errors.Unwrap()
// on errors of that type always returns nil as that function cannot
// return multiple errors.
//
// The error type is not exported so can only used through the standard
// Go error interfaces.
package errs

import (
	"errors"
	"fmt"
	"strings"
)

// New takes a variable number of input errors and returns an error that
// is formatted as the concatenation of the string form of each input
// error, separated by a colon and space. The returned error wraps each
// input error.
//
// For example New(err1, err2, err3) is formatted as
//    err1: err2: err3
func New(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	format := strings.Repeat("%v: ", len(errs)-1) + "%v"
	args := make([]interface{}, len(errs))
	for i := range errs {
		args[i] = errs[i]
	}
	return Errorf(format, args...)
}

// Errorf returns an error formatted with a format string and arguments,
// similar to how fmt.Errorf works. The error returned by Errorf wraps
// each argument that is an error type unless that error argument is
// marked with the Nowrap() function. The %w format verb should not be
// used in the format string.
func Errorf(format string, args ...interface{}) error {
	rawArgs := make([]interface{}, len(args))
	var errs []error
	for i, arg := range args {
		if n, ok := arg.(noWrap); ok {
			arg = n.err
		} else if e, ok := arg.(error); ok {
			errs = append(errs, e)
		}
		rawArgs[i] = arg
	}
	return &multiErr{
		s:    fmt.Sprintf(format, rawArgs...),
		errs: errs,
	}
}

// NoWrap marks errors not to be wrapped for Errorf.
func NoWrap(err error) noWrap { //nolint:golint
	return noWrap{err: err}
}

type noWrap struct{ err error }

// multiErr contains wrapped errors and a formatted error message.
type multiErr struct {
	s    string
	errs []error
}

// Error returns multiErr's error message to implement error interface.
func (e *multiErr) Error() string {
	return e.s
}

// Is returns true if any error in multiErr's slice of errors or any
// error within each errors chain matches the target value, or false if
// not.
func (e *multiErr) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// As finds the first error in multiErr's slice of errors or any error
// within each error's chain that matches the target type. If such an
// error is found, target is set to the error value and true is
// returned. Otherwise false is returned.
func (e *multiErr) As(target interface{}) bool {
	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

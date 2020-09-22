// Package httpe provides a model of HTTP handler that returns errors.
//
// The standard Go library package net/http provides a Handler interface that
// requires the handler write errors to the provided ResponseWriter. This is
// different to the usual Go way of handling errors that has functions
// returning errors, and it makes normal http.Handlers a bit cumbersome and
// repetitive in the error handling cases.
//
// This package provides a HandlerE interface with a ServeHTTPe method that has
// the same signature as http.Handler.ServeHTTP except it also returns an
// error. A separate error handler can be bound to the HandlerE using
// NewHandler() and turn the HandlerE into an http.Handler.
//
// As well as making handler code a little simpler, separating the ErrWriter
// allows for common error handling amongst disparate handlers.
//
// The default ErrWriter writes errors that wrap an httpe.StatusError by
// writing the status code of that error and if the status code is a client
// error, writes the error formatted as a string. If it is not a client error,
// it writes just the text for that status code. Errors that do not wrap an
// httpe.StatusError are treated as httpe.ErrInternalServerError.
//
// Option arguments to NewHandler() allow a custom ErrWriter to be provided.
package httpe

import (
	"fmt"
	"net/http"
)

// HandlerE works like an HTTP.Handler with the addition of an error return
// value. It is intended to be used with ErrWriter which handles the error and
// writes an appropriate http response StatusCode and Body to the
// ResponseWriter.
type HandlerE interface {
	ServeHTTPe(http.ResponseWriter, *http.Request) error
}

// The HandlerFuncE type is an adapter to allow the use of ordinary
// functions as HandlerE. If f is a function with the appropriate
// signature, HandlerFuncE(f) is a HandlerE that calls f.
type HandlerFuncE func(http.ResponseWriter, *http.Request) error

// ServeHTTPe calls f(w, r) and returns its error.
func (f HandlerFuncE) ServeHTTPe(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// ErrWriter translates an error into the appropriate http response
// StatusCode and Body and writes it.
type ErrWriter interface {
	WriteErr(http.ResponseWriter, error)
}

// The ErrWriterFunc type is an adapter to allow the use of ordinary
// functions as ErrWriter. If f is a function with the appropriate
// signature, ErrWriterFunc(f) is a ErrWriter that calls f.
type ErrWriterFunc func(http.ResponseWriter, error)

// WriteErr translates an error into the appropriate http response
// StatusCode and Body and writes it.
func (ew ErrWriterFunc) WriteErr(w http.ResponseWriter, err error) {
	ew(w, err)
}

// New returns an http.Handler that calls in sequence all the args that are a
// type of handler, stopping if any return an error. If any of the args is an
// ErrWriter or a function that has the signature of an ErrWriterFunc, it will
// be called to handle the error if there was one.
//
// The types that are recognised as handlers in the arg list are any type that
// implements HandlerE (including HandlerFuncE), a function that matches the
// signature of a HandlerFuncE, an http.Handler, or a function that matches the
// signature of an http.HandlerFunc. Args of the latter two are adapted to
// always return a nil error.
//
// If an argument does not match any of the preceding types or more than one
// ErrWriter is passed, an error is returned.
func New(args ...interface{}) (http.Handler, error) {
	handlers := make([]HandlerE, 0, len(args))
	opts := []option{}

	for i, arg := range args {
		switch v := arg.(type) {
		case HandlerE:
			handlers = append(handlers, v)
		case func(http.ResponseWriter, *http.Request) error:
			handlers = append(handlers, HandlerFuncE(v))
		case http.Handler:
			handlers = append(handlers, handlerAdapter(v))
		case func(http.ResponseWriter, *http.Request):
			handlers = append(handlers, handlerAdapter(http.HandlerFunc(v)))
		case ErrWriter:
			opts = append(opts, WithErrWriter(v))
		case func(http.ResponseWriter, error):
			opts = append(opts, WithErrWriterFunc(v))
		default:
			return nil, fmt.Errorf("arg %d: unknown arg type: %T", i, v)
		}
		if len(opts) > 1 {
			return nil, fmt.Errorf("arg %d: too many ErrWriters", i)
		}
	}
	return NewHandler(Chain(handlers...), opts...), nil
}

// Must passes all its args to New() and panics if New() returns an error. If
// it does not, the handler result of New() is returned.
func Must(args ...interface{}) http.Handler {
	h, err := New(args...)
	if err != nil {
		panic(err)
	}
	return h
}

// handlerAdapter turns an http.Handler into a HandlerE that returns nil.
func handlerAdapter(h http.Handler) HandlerE {
	f := func(w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTP(w, r)
		return nil
	}
	return HandlerFuncE(f)
}

// NewHandler returns an http.Handler that calls h.ServeHTTPe and handles the
// error returned, if any, with an ErrWriter to write the error to the
// ResponseWriter. The default ErrWriter is httpe.WriteSafeErr but can be
// overridden with an option passed to NewHandler.
func NewHandler(h HandlerE, opts ...option) http.Handler {
	o := newOptions(opts)
	f := func(w http.ResponseWriter, r *http.Request) {
		if err := h.ServeHTTPe(w, r); err != nil {
			o.ew.WriteErr(w, err)
		}
	}
	return http.HandlerFunc(f)
}

// NewHandlerFunc returns an http.HandlerFunc that calls h() and handles the
// error returned, if any, with an ErrWriter to write the error to the
// ResponseWriter. The default ErrWriter is httpe.WriteSafeErr but can be
// overridden with an option passed to NewHandlerFunc.
func NewHandlerFunc(h HandlerFuncE, opts ...option) http.HandlerFunc { //nolint:interfacer
	return NewHandler(h, opts...).(http.HandlerFunc)
}

type option func(*options)

type options struct {
	ew ErrWriter
}

func newOptions(opts []option) options {
	o := options{
		ew: ErrWriterFunc(WriteSafeErr),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithErrWriter returns an option to use the given ErrWriter as the ErrWriter
// for a HandlerE.
func WithErrWriter(ew ErrWriter) option { //nolint:golint // Do not want to export option type.
	return func(o *options) {
		o.ew = ew
	}
}

// WithErrWriterFunc returns an option to use the given ErrWriterFunc as the
// ErrWriter for a HandlerE.
func WithErrWriterFunc(f ErrWriterFunc) option { //nolint:golint // Do not want to export option type.
	return func(o *options) {
		o.ew = f
	}
}

// Chain returns a HandlerE that executes each of the HandlerFuncE parameters
// sequentially, stopping at the first one that returns an error and returning
// that error. It returns nil if none of the handlers return an error.
func Chain(he ...HandlerE) HandlerE {
	f := func(w http.ResponseWriter, r *http.Request) error {
		for _, h := range he {
			if err := h.ServeHTTPe(w, r); err != nil {
				return err
			}
		}
		return nil
	}
	return HandlerFuncE(f)
}

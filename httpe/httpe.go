// Package httpe provides a model of HTTP handler that returns errors.
//
// net/http provides a Handler interface that requires the handler to write
// errors to the provided ResponseWriter. This is different to the usual Go way
// of handling errors that has functions returning errors, and it makes normal
// http.Handlers a bit cumbersome.
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

// Chain returns a HandlerE which executes each of the HandlerFuncE
// parameters sequentially stopping at the first one that returns an
// error, and returning that error, or nil if none return an error.
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

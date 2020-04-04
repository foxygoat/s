// Package httpe holds http server utilities
package httpe

import (
	"net/http"
)

// HandlerE ...
type HandlerE interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// HandlerFuncE works like an HTTP.HandlerFunc with the addition of an error
// return value. It is intended to be used with WriteErr which handles
// the error and turns it into an appropriate http response StatusCode and
// Body.
type HandlerFuncE func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r) and returns its error.
func (f HandlerFuncE) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// ErrWriter translates an error into the appropriate http response
// StatusCode and Body and writes it.
type ErrWriter interface {
	Write(http.ResponseWriter, error)
}

// ErrWriterFunc type is an adapter to allow the use of ordinary
// functions as ErrWriter. If f is a function with the appropriate
// signature, ErrWriterFunc(f) is a ErrWriter that calls f.
type ErrWriterFunc func(http.ResponseWriter, error)

func (ew ErrWriterFunc) Write(w http.ResponseWriter, err error) {
	ew(w, err)
}

// NewHandler creates a new http.HandlerFunc which calls
// HandlerE.ServeHTTP and if it returns an error calls WriteErr to
// create the appropriate response.
func NewHandler(h HandlerE, ew ErrWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.ServeHTTP(w, r); err != nil {
			ew.Write(w, err)
		}
	}
}

// Chain returns a HandlerFuncE which executes each of the HandlerFuncE
// parameters sequentially stopping at the first one that returns an
// error, and returning that error, or nil if none return an error.
func Chain(hf ...HandlerE) HandlerFuncE {
	return func(w http.ResponseWriter, r *http.Request) error {
		for _, h := range hf {
			if err := h.ServeHTTP(w, r); err != nil {
				return err
			}
		}
		return nil
	}
}

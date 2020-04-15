// Package httpe holds http server utilities
package httpe

import (
	"net/http"
)

// HandlerE works like an HTTP.Handler with the addition of an error
// return value. It is intended to be used with ErrWriter which handles
// the error and turns it into an appropriate http response StatusCode
// and Body.
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

// NewHandler creates a new http.HandlerFunc which calls
// HandlerE.ServeHTTP and if it returns an error calls ErrWriter.Write
// to create the appropriate response.
func NewHandler(h HandlerE, ew ErrWriter) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if err := h.ServeHTTPe(w, r); err != nil {
			ew.WriteErr(w, err)
		}
	}
	return http.HandlerFunc(f)
}

// NewHandlerFunc creates a new http.HandlerFunc which calls
// HandlerFuncE and if it returns an error calls ErrWriterFunc to
// create the appropriate response.
func NewHandlerFunc(h HandlerFuncE, ew ErrWriterFunc) http.HandlerFunc { //nolint:interfacer
	return NewHandler(h, ew).(http.HandlerFunc)
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

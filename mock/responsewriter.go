// Package mock provides mocking utilities used in tests.
package mock

import (
	"net/http"
)

// ResponseWriter returns a responseWriter implementing
// http.ResponseWriter.
func ResponseWriter() *responseWriter { //nolint:golint
	return &responseWriter{}
}

// responseWriter implements http.ResponseWriter.
type responseWriter struct {
	LastBody   string
	LastStatus int
	header     http.Header
	err        error
}

// Err sets the error value returned by the Write method.
// Chain it with construction: mock.ResponseWriter().Err(someErr).
func (r *responseWriter) Err(err error) *responseWriter {
	r.err = err
	return r
}

// Header returns the responseWriters HTTP header.
func (r *responseWriter) Header() http.Header {
	if r.header == nil {
		r.header = http.Header{}
	}
	return r.header
}

// WriteHeader writes HTTP Status code to struct field.
func (r *responseWriter) WriteHeader(status int) {
	r.LastStatus = status
}

// Write writes response body to struct field.
func (r *responseWriter) Write(b []byte) (int, error) {
	r.LastBody = string(b)
	return len(b), r.err
}

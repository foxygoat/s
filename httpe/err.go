package httpe

import (
	"errors"
	"net/http"
)

// StatusError wraps an http.StatusCode of value 4xx or 5xx. It is used
// in combination with various sentinel values each representing a http
// status code for convenience with HandlerE and HandlerFuncE.
type StatusError int

// Error returns the error message and implements the error interface.
func (err StatusError) Error() string { return http.StatusText(int(err)) }

// IsClientError returns true if the error is in range 400 to 499.
func (err StatusError) IsClientError() bool { return int(err) >= 400 && int(err) <= 499 }

// Code returns a status int representing an http.Status* value.
func (err StatusError) Code() int { return int(err) }

// HTTP status codes as registered with IANA.
// See: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
var (
	ErrBadRequest                   = StatusError(http.StatusBadRequest)
	ErrUnauthorized                 = StatusError(http.StatusUnauthorized)
	ErrPaymentRequired              = StatusError(http.StatusPaymentRequired)
	ErrForbidden                    = StatusError(http.StatusForbidden)
	ErrNotFound                     = StatusError(http.StatusNotFound)
	ErrMethodNotAllowed             = StatusError(http.StatusMethodNotAllowed)
	ErrNotAcceptable                = StatusError(http.StatusNotAcceptable)
	ErrProxyAuthRequired            = StatusError(http.StatusProxyAuthRequired)
	ErrRequestTimeout               = StatusError(http.StatusRequestTimeout)
	ErrConflict                     = StatusError(http.StatusConflict)
	ErrGone                         = StatusError(http.StatusGone)
	ErrLengthRequired               = StatusError(http.StatusLengthRequired)
	ErrPreconditionFailed           = StatusError(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge        = StatusError(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong            = StatusError(http.StatusRequestURITooLong)
	ErrUnsupportedMediaType         = StatusError(http.StatusUnsupportedMediaType)
	ErrRequestedRangeNotSatisfiable = StatusError(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed            = StatusError(http.StatusExpectationFailed)
	ErrTeapot                       = StatusError(http.StatusTeapot)
	ErrMisdirectedRequest           = StatusError(http.StatusMisdirectedRequest)
	ErrUnprocessableEntity          = StatusError(http.StatusUnprocessableEntity)
	ErrLocked                       = StatusError(http.StatusLocked)
	ErrFailedDependency             = StatusError(http.StatusFailedDependency)
	ErrTooEarly                     = StatusError(http.StatusTooEarly)
	ErrUpgradeRequired              = StatusError(http.StatusUpgradeRequired)
	ErrPreconditionRequired         = StatusError(http.StatusPreconditionRequired)
	ErrTooManyRequests              = StatusError(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge  = StatusError(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons   = StatusError(http.StatusUnavailableForLegalReasons)

	ErrInternalServerError           = StatusError(http.StatusInternalServerError)
	ErrNotImplemented                = StatusError(http.StatusNotImplemented)
	ErrBadGateway                    = StatusError(http.StatusBadGateway)
	ErrServiceUnavailable            = StatusError(http.StatusServiceUnavailable)
	ErrGatewayTimeout                = StatusError(http.StatusGatewayTimeout)
	ErrHTTPVersionNotSupported       = StatusError(http.StatusHTTPVersionNotSupported)
	ErrVariantAlsoNegotiates         = StatusError(http.StatusVariantAlsoNegotiates)
	ErrInsufficientStorage           = StatusError(http.StatusInsufficientStorage)
	ErrLoopDetected                  = StatusError(http.StatusLoopDetected)
	ErrNotExtended                   = StatusError(http.StatusNotExtended)
	ErrNetworkAuthenticationRequired = StatusError(http.StatusNetworkAuthenticationRequired)
)

// WriteSafeErr writes err as an HTTP error to the http.ResponseWriter.
//
// If err wraps a StatusError, WriteSafeErr writes the Code of that error as
// the HTTP status code. Otherwise writes 500 as the code (internal server
// error).
//
// WriteSafeErr writes the HTTP response body as the error described as a
// string if the error wraps a StatusError and the Code for that error is a
// client error (code 400 to 499). Otherwise, only the text for the status code
// is written. Errors that do not wrap a StatusError are treated as
// ErrInternalServerError. This prevents leaking internal details from a server
// error into the response to the client, but allows adding information to a
// client error to inform the client of details of that error.
func WriteSafeErr(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	sErr := ErrInternalServerError
	if ok := errors.As(err, &sErr); !ok || !sErr.IsClientError() {
		// Hide the actual error to prevent information leakage
		err = sErr
	}
	http.Error(w, err.Error(), sErr.Code())
}

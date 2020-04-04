package httpe

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	// ErrInvalidHeader ...
	ErrInvalidHeader = fmt.Errorf("invalid header")
	// ErrMalformedJSON ...
	ErrMalformedJSON = fmt.Errorf("malformed JSON")
	// ErrMethodNotAllowed ...
	ErrMethodNotAllowed = Err(nil, http.StatusMethodNotAllowed)
)

// Error ...
type Error struct {
	Err        error
	HTTPStatus int
}

// Err ...
func Err(err error, httpStatus int) error {
	return &Error{Err: err, HTTPStatus: httpStatus}
}

func wrapError(err error) *Error {
	var e *Error
	if ok := errors.As(err, &e); ok {
		return e
	}
	return &Error{Err: err, HTTPStatus: http.StatusInternalServerError}
}

// Error ...
func (e *Error) Error() string {
	if e.Err == nil {
		return strings.ToLower(http.StatusText(e.HTTPStatus))
	}
	return e.Err.Error()
}

// HTTPMessage ...
func (e *Error) HTTPMessage() string {
	if e.Err == nil || e.HTTPStatus == 0 || e.HTTPStatus == http.StatusInternalServerError {
		return "internal server error"
	}
	return e.Error()
}

// WriteErrTxt ...
func WriteErrTxt(w http.ResponseWriter, err error) {
	e := wrapError(err)
	http.Error(w, e.HTTPMessage(), e.HTTPStatus)
}

// WriteErrJSON ...
func WriteErrJSON(w http.ResponseWriter, err error) {
	e := wrapError(err)
	b, _ := json.Marshal(map[string]string{
		"errorMessage": e.HTTPMessage(),
		"code":         strconv.Itoa(e.HTTPStatus)})
	http.Error(w, string(b), e.HTTPStatus)
}

package httpe

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"foxygo.at/s/mock"
	"github.com/stretchr/testify/require"
)

var errNoPost = fmt.Errorf("no POST method")

type handlerE struct{}

func (*handlerE) ServeHTTPe(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errNoPost
	}
	return nil
}

func (*handlerE) WriteErr(w http.ResponseWriter, err error) {
	fmt.Fprint(w, err.Error())
}

func TestNew(t *testing.T) {
	var handlerE HandlerE = HandlerFuncE(func(w http.ResponseWriter, _ *http.Request) error {
		fmt.Fprintf(w, "1")
		return nil
	})
	funcE := func(w http.ResponseWriter, _ *http.Request) error {
		fmt.Fprintf(w, "2")
		return nil
	}
	var handlerH http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "3")
	})
	funcH := func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "4")
	}

	h, err := New(handlerE, funcE, handlerH, funcH)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, &http.Request{})
	require.Equal(t, "1234", w.Body.String())

	f := func() { Must(handlerE, funcE, handlerH, funcH) }
	require.NotPanics(t, f)
}

func TestNewErr(t *testing.T) {
	ew := ErrWriterFunc(func(_ http.ResponseWriter, _ error) {})

	// Test multuple ErrWriters
	_, err := New(ew, func(_ http.ResponseWriter, _ error) {})
	require.Error(t, err)

	// Test unknown type
	_, err = New(42)
	require.Error(t, err)

	f := func() { Must(ew, ew) }
	require.Panics(t, f)
}

func TestNewErrWriter(t *testing.T) {
	funcE := func(w http.ResponseWriter, _ *http.Request) error {
		return fmt.Errorf("%w: 🐿️", ErrBadRequest)
	}
	ew := ErrWriterFunc(func(w http.ResponseWriter, err error) {
		fmt.Fprint(w, "error: ", err.Error())
	})

	h, err := New(funcE, ew)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, &http.Request{})
	require.Equal(t, "error: Bad Request: 🐿️", w.Body.String())
}

func TestHandler(t *testing.T) {
	he := &handlerE{}
	h := NewHandler(he)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.LastStatus)
}

func TestHandlerFunc(t *testing.T) {
	he := handlerE{}
	h := NewHandlerFunc(he.ServeHTTPe)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.LastStatus)
}

func TestHandlerWithErrWriter(t *testing.T) {
	he := &handlerE{}
	h := NewHandler(he, WithErrWriter(he))
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, errNoPost.Error(), w.LastBody)
}

func TestHandlerFuncWithErrWriter(t *testing.T) {
	he := handlerE{}
	h := NewHandlerFunc(he.ServeHTTPe, WithErrWriterFunc(he.WriteErr))
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, errNoPost.Error(), w.LastBody)
}

func TestChain(t *testing.T) {
	count := 0
	errHand := fmt.Errorf("error 🤚")
	goodHandler := HandlerFuncE(func(http.ResponseWriter, *http.Request) error {
		count++
		return nil
	})
	errorHandler := HandlerFuncE(func(http.ResponseWriter, *http.Request) error {
		return errHand
	})

	h := Chain(goodHandler, goodHandler, goodHandler)
	err := h.ServeHTTPe(mock.ResponseWriter(), &http.Request{})
	require.NoError(t, err)
	require.Equal(t, count, 3)

	count = 0
	h = Chain(goodHandler, errorHandler, goodHandler)
	err = h.ServeHTTPe(mock.ResponseWriter(), &http.Request{})
	require.Error(t, err)
	require.Equal(t, err, errHand)
	require.Equal(t, count, 1)
}

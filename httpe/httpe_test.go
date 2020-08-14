package httpe

import (
	"fmt"
	"net/http"
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

func TestHandler(t *testing.T) {
	he := &handlerE{}
	h := NewHandler(he, he)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, errNoPost.Error(), w.LastBody)
}

func TestHandlerFunc(t *testing.T) {
	he := handlerE{}
	h := NewHandlerFunc(he.ServeHTTPe, he.WriteErr)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, errNoPost.Error(), w.LastBody)
}

func TestSafeHandler(t *testing.T) {
	he := &handlerE{}
	h := NewSafeHandler(he)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.LastStatus)
}

func TestSafeHandlerFunc(t *testing.T) {
	he := handlerE{}
	h := NewSafeHandlerFunc(he.ServeHTTPe)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.LastStatus)
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

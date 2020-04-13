package httpe

import (
	"fmt"
	"net/http"
	"testing"

	"foxygo.at/s/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	errWriter := ErrWriterFunc(func(w http.ResponseWriter, err error) {
		fmt.Fprint(w, err.Error())
	})

	errNoPost := fmt.Errorf("no POST method")
	handlerE := HandlerFuncE(func(_ http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return errNoPost
		}
		return nil
	})
	h := NewHandler(handlerE, errWriter)
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
	chain := Chain(goodHandler, goodHandler, goodHandler)
	err := chain(mock.ResponseWriter(), &http.Request{})
	require.NoError(t, err)
	require.Equal(t, count, 3)
	count = 0
	chain = Chain(goodHandler, errorHandler, goodHandler)
	err = chain(mock.ResponseWriter(), &http.Request{})
	require.Error(t, err)
	require.Equal(t, err, errHand)
	require.Equal(t, count, 1)
}

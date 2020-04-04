package httpe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"foxygo.at/s/mock"
	"github.com/stretchr/testify/require"
)

func TestHeaderCheckHandler(t *testing.T) {
	requiredHeader := http.Header{"Content-Type": []string{"application/json"}}
	h := NewHeaderCheck(requiredHeader)
	req := &http.Request{Header: requiredHeader}
	require.NoError(t, h.ServeHTTP(nil, req))
	req = &http.Request{Header: nil}
	err := h.ServeHTTP(nil, req)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidHeader))
}

func TestCheckPost(t *testing.T) {
	req := &http.Request{Method: http.MethodPost}
	require.NoError(t, CheckPost(nil, req))
	req.Method = http.MethodGet
	err := CheckPost(nil, req)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrMethodNotAllowed))
}

func TestHandler(t *testing.T) {
	ew := ErrWriterFunc(func(w http.ResponseWriter, err error) {
		_, _ = w.Write([]byte(err.Error()))
	})

	h := NewHandler(CheckPost, ew)
	r := &http.Request{Method: http.MethodGet}
	w := mock.ResponseWriter()
	h.ServeHTTP(w, r)
	require.Equal(t, ErrMethodNotAllowed.Error(), w.LastBody)
}

func TestChain(t *testing.T) {
	count := 0
	testErr := fmt.Errorf("error")
	goodHandler := HandlerFuncE(func(http.ResponseWriter, *http.Request) error {
		count++
		return nil
	})
	errorHandler := HandlerFuncE(func(http.ResponseWriter, *http.Request) error {
		return testErr
	})
	chain := Chain(goodHandler, goodHandler, goodHandler)
	err := chain(mock.ResponseWriter(), &http.Request{})
	require.NoError(t, err)
	require.Equal(t, count, 3)
	count = 0
	chain = Chain(goodHandler, errorHandler, goodHandler)
	err = chain(mock.ResponseWriter(), &http.Request{})
	require.Error(t, err)
	require.Equal(t, err, testErr)
	require.Equal(t, count, 1)
}

func TestWriteJSON(t *testing.T) {
	w := mock.ResponseWriter()
	err := WriteJSON(w, "validJSON")
	require.NoError(t, err)
	require.Equal(t, `"validJSON"`+"\n", w.LastBody)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWriteJSONMarshalErr(t *testing.T) {
	w := mock.ResponseWriter()
	ch := make(chan (int))
	err := WriteJSON(w, ch)
	require.Error(t, err)
	jsonErr := &json.UnsupportedTypeError{}
	require.True(t, errors.As(err, &jsonErr))
}

func TestWriteJSONWriterErr(t *testing.T) {
	errMock := fmt.Errorf("mock writer error")
	w := mock.ResponseWriter().Err(errMock)
	err := WriteJSON(w, "validJSON")
	require.Error(t, err)
	require.Equal(t, errMock, err)
}

func TestReadJSON(t *testing.T) {
	b := bytes.NewBufferString(`"validJSON"`)
	var target string
	err := ReadJSON(b, &target)
	require.NoError(t, err)
	require.Equal(t, "validJSON", target)
}

func TestReadJSONErr(t *testing.T) {
	b := bytes.NewBufferString(`{"invalidJSON":`)
	err := ReadJSON(b, &struct{}{})
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrMalformedJSON))
	b = bytes.NewBufferString(`"validJSON"`)
	err = ReadJSON(b, &struct{}{})
	require.Error(t, err)
	require.False(t, errors.Is(err, ErrMalformedJSON))
	unmarshalTypeErr := &json.UnmarshalTypeError{}
	require.True(t, errors.As(err, &unmarshalTypeErr))
}

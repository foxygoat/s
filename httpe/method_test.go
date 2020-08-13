package httpe

import (
	"errors"
	"net/http"
	"testing"

	"foxygo.at/s/mock"
	"github.com/stretchr/testify/require"
)

func TestEnsureMethod(t *testing.T) {
	tests := map[string]HandlerFuncE{
		http.MethodGet:     Get,
		http.MethodHead:    Head,
		http.MethodPost:    Post,
		http.MethodPut:     Put,
		http.MethodPatch:   Patch,
		http.MethodDelete:  Delete,
		http.MethodConnect: Connect,
		http.MethodOptions: Options,
		http.MethodTrace:   Trace,
	}
	for method, h := range tests {
		r := &http.Request{Method: method}
		w := mock.ResponseWriter()
		err := h.ServeHTTPe(w, r)
		require.NoError(t, err)
	}

	r := &http.Request{Method: http.MethodPost}
	w := mock.ResponseWriter()
	err := Get.ServeHTTPe(w, r)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrMethodNotAllowed))
}

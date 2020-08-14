package httpe

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"foxygo.at/s/mock"
	"github.com/stretchr/testify/require"
)

func TestWriteSafeErr(t *testing.T) {
	errs := []error{
		ErrBadRequest, ErrUnauthorized,
		ErrPaymentRequired, ErrForbidden, ErrNotFound,
		ErrMethodNotAllowed, ErrNotAcceptable,
		ErrProxyAuthRequired, ErrRequestTimeout,
		ErrConflict, ErrGone, ErrLengthRequired,
		ErrPreconditionFailed, ErrRequestEntityTooLarge,
		ErrRequestURITooLong, ErrUnsupportedMediaType,
		ErrRequestedRangeNotSatisfiable, ErrExpectationFailed,
		ErrTeapot, ErrMisdirectedRequest,
		ErrUnprocessableEntity, ErrLocked,
		ErrFailedDependency, ErrTooEarly,
		ErrUpgradeRequired, ErrPreconditionRequired,
		ErrTooManyRequests, ErrRequestHeaderFieldsTooLarge,
		ErrUnavailableForLegalReasons, ErrInternalServerError,
		ErrNotImplemented, ErrBadGateway,
		ErrServiceUnavailable, ErrGatewayTimeout,
		ErrHTTPVersionNotSupported, ErrVariantAlsoNegotiates,
		ErrInsufficientStorage, ErrLoopDetected,
		ErrNotExtended, ErrNetworkAuthenticationRequired,
	}
	for _, err := range errs {
		w := httptest.NewRecorder()
		WriteSafeErr(w, err)
		require.Contains(t, err.Error(), http.StatusText(w.Code))
	}
}

func TestInternalServerErrDefault(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("💥")
	WriteSafeErr(w, err)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestWrappedErr(t *testing.T) {
	w := httptest.NewRecorder()
	err2 := fmt.Errorf("%w: feeling very lost", ErrNotFound)
	WriteSafeErr(w, err2)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestNoErr(t *testing.T) {
	w := httptest.NewRecorder()
	WriteSafeErr(w, nil)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestNoLeak(t *testing.T) {
	w := mock.ResponseWriter()
	err := fmt.Errorf("%w: secret details", ErrInternalServerError)
	WriteSafeErr(w, err)
	require.NotContains(t, w.LastBody, "secret")

	err = fmt.Errorf("%w: secret details", ErrBadRequest)
	WriteSafeErr(w, err)
	require.Contains(t, w.LastBody, "secret")
}

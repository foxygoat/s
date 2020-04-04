package httpe

import (
	"fmt"
	"net/http"
)

// CheckJSONHeader ...
var CheckJSONHeader = NewHeaderCheck(http.Header{"Content-Type": []string{"application/json"}})

// NewHeaderCheck returns a HandlerFuncE that ensures that all required
// headers are present in the request headers. Additional headers are
// allowed.
func NewHeaderCheck(requiredHeader http.Header) HandlerFuncE {
	return func(_ http.ResponseWriter, r *http.Request) error {
		return ValidateHeader(requiredHeader, r.Header)
	}
}

// ValidateHeader ensures that got Headers are equal or a super set of
// wanted headers. Wanted headers need to be canonical, got headers
// are canonicalised. If a header with given key is provided in
// wanted Headers it needs to be present in the same way in go headers.
func ValidateHeader(want, got http.Header) error {
	canonicalGot := http.Header{}
	for k, v := range got {
		canonicalGot[http.CanonicalHeaderKey(k)] = v
	}
	for k, v := range want {
		if !isSubset(v, canonicalGot[k]) {
			return fmt.Errorf("%w: want %#v, got %#v", ErrInvalidHeader, v, canonicalGot[k])
		}
	}
	return nil
}

func isSubset(sub, slice []string) bool {
	for _, str := range sub {
		if !contains(slice, str) {
			return false
		}
	}
	return true
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

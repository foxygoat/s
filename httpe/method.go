package httpe

import (
	"net/http"
)

// CheckGet ...
var CheckGet = NewMethodCheck(http.MethodGet)

// CheckPost ...
var CheckPost = NewMethodCheck(http.MethodPost)

// CheckPut ...
var CheckPut = NewMethodCheck(http.MethodPut)

// CheckPatch ...
var CheckPatch = NewMethodCheck(http.MethodPatch)

// CheckDelete ...
var CheckDelete = NewMethodCheck(http.MethodDelete)

// NewMethodCheck ...
func NewMethodCheck(method string) HandlerFuncE {
	return HandlerFuncE(func(_ http.ResponseWriter, r *http.Request) error {
		if r.Method != method {
			return ErrMethodNotAllowed
		}
		return nil
	})
}

// MethodMux .... string: http.Method
type MethodMux map[string]HandlerE

// ServeHTTP ...
func (m MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	h := m[r.Method]
	if h == nil {
		return ErrMethodNotAllowed
	}
	return h.ServeHTTP(w, r)
}

// NewMethodMux ...
func NewMethodMux() MethodMux {
	return MethodMux{}
}

// Get ...
func (m MethodMux) Get(h HandlerE) MethodMux {
	m[http.MethodGet] = h
	return m
}

// Post ....
func (m MethodMux) Post(h HandlerE) MethodMux {
	m[http.MethodPost] = h
	return m
}

// Put ...
func (m MethodMux) Put(h HandlerE) MethodMux {
	m[http.MethodPut] = h
	return m
}

// Delete ...
func (m MethodMux) Delete(h HandlerE) MethodMux {
	m[http.MethodDelete] = h
	return m
}

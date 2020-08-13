package httpe

import "net/http"

var (
	// Get returns a ErrMethodNotAllowed if request method is not a GET. Use with Chain.
	Get = newMethodChecker(http.MethodGet)
	// Head returns a ErrMethodNotAllowed if request method is not a HEAD. Use with Chain.
	Head = newMethodChecker(http.MethodHead)
	// Post returns a ErrMethodNotAllowed if request method is not a Post. Use with Chain.
	Post = newMethodChecker(http.MethodPost)
	// Put returns a ErrMethodNotAllowed if request method is not a Put. Use with Chain.
	Put = newMethodChecker(http.MethodPut)
	// Patch returns a ErrMethodNotAllowed if request method is not a Patch. Use with Chain.
	Patch = newMethodChecker(http.MethodPatch)
	// Delete returns a ErrMethodNotAllowed if request method is not a Delete. Use with Chain.
	Delete = newMethodChecker(http.MethodDelete)
	// Connect returns a ErrMethodNotAllowed if request method is not a Connect. Use with Chain.
	Connect = newMethodChecker(http.MethodConnect)
	// Options returns a ErrMethodNotAllowed if request method is not a Options. Use with Chain.
	Options = newMethodChecker(http.MethodOptions)
	// Trace returns a ErrMethodNotAllowed if request method is not a Trace. Use with Chain.
	Trace = newMethodChecker(http.MethodTrace)
)

func newMethodChecker(method string) HandlerFuncE {
	return HandlerFuncE(func(_ http.ResponseWriter, r *http.Request) error {
		if r.Method != method {
			return ErrMethodNotAllowed
		}
		return nil
	})
}

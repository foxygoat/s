package httpe

import "net/http"

var (
	// Get is a HandlerE that returns a ErrMethodNotAllowed if the request
	// method is not GET. Use with Chain or New/Must.
	Get = newMethodChecker(http.MethodGet)

	// Head is a HandlerE that returns a ErrMethodNotAllowed if the request
	// method is not HEAD. Use with Chain or New/Must.
	Head = newMethodChecker(http.MethodHead)

	// Post is a HandlerE that returns a ErrMethodNotAllowed if the request
	// method is not POST. Use with Chain or New/Must.
	Post = newMethodChecker(http.MethodPost)

	// Put is a HandlerE that returns a ErrMethodNotAllowed if the request
	// method is not PUT. Use with Chain or New/Must.
	Put = newMethodChecker(http.MethodPut)

	// Patch is a HandlerE that returns a ErrMethodNotAllowed if the
	// request method is not PATCH. Use with Chain or New/Must.
	Patch = newMethodChecker(http.MethodPatch)

	// Delete is a HandlerE that returns a ErrMethodNotAllowed if the
	// request method is not DELETE. Use with Chain or New/Must.
	Delete = newMethodChecker(http.MethodDelete)

	// Connect is a HandlerE that returns a ErrMethodNotAllowed if the
	// request method is not CONNECT. Use with Chain or New/Must.
	Connect = newMethodChecker(http.MethodConnect)

	// Options is a HandlerE that returns a ErrMethodNotAllowed if the
	// request method is not OPTIONS. Use with Chain or New/Must.
	Options = newMethodChecker(http.MethodOptions)

	// Trace is a HandlerE that returns a ErrMethodNotAllowed if the
	// request method is not TRACE. Use with Chain or New/Must.
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

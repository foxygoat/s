package httpe_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"foxygo.at/s/httpe"
)

func ExampleNew() {
	// Create a handler by chaining a number of other handlers and an
	// error writer to handle any errors from those handlers.
	handler, _ := httpe.New(corsAllowAll, httpe.Get, &api{}, errWriter)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)
	handler.ServeHTTP(w, r)
	fmt.Printf("%d %s\n", w.Code, w.Body.String()) // 200 world

	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/hello", nil)
	handler.ServeHTTP(w, r)
	fmt.Printf("%d %s\n", w.Code, w.Body.String()) // 405 🐈 METHOD NOT ALLOWED!!!1!

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/goodbye", nil)
	handler.ServeHTTP(w, r)
	fmt.Printf("%d %s\n", w.Code, w.Body.String()) // 500 🐈 NO GOODBYES!!!1!

	// output:
	// 200 world
	// 405 🐈 METHOD NOT ALLOWED!!!1!
	// 500 🐈 NO GOODBYES!!!1!
}

// A traditional http.Handler compatible function.
func corsAllowAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

type api struct{}

// Implements httpe.HandlerE.
func (a *api) ServeHTTPe(w http.ResponseWriter, r *http.Request) error {
	switch r.URL.Path {
	case "/hello":
		fmt.Fprintf(w, "world")
	case "/goodbye":
		return errors.New("no goodbyes")
	}
	return nil
}

// Matches httpe.ErrWriterFunc.
func errWriter(w http.ResponseWriter, err error) {
	var se httpe.StatusError
	if errors.As(err, &se) {
		w.WriteHeader(se.Code())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "🐈 "+strings.ToUpper(err.Error())+"!!!1!")
}

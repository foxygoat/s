package httpe_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"foxygo.at/s/httpe"
)

type User struct {
	Name string
	Age  int
}

var (
	errInput     = errors.New("input error")
	errDuplicate = errors.New("duplicate")
	storage      = map[string]User{}
)

func ExampleHandlerFuncE() {
	// Create a http.HandlerFunc from our httpe.HandlerFuncE function and
	// an httpe.ErrWriterFunc function.
	handler := httpe.NewHandlerFunc(handle, httpe.WithErrWriterFunc(writeErr))

	// In this example, we call the handler directly using httptest.
	// Normally you would start a http server.
	// http.ListenAndServe(":9090", handler)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/user", strings.NewReader(`{"Name": "truncated...`))
	handler(w, r)
	fmt.Printf("%d %s", w.Code, w.Body.String())
	// output: 400 input error
}

func handle(w http.ResponseWriter, r *http.Request) error {
	user := User{}
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &user); err != nil {
		return errInput
	}
	if user.Name == "" || user.Age < 0 {
		return errInput
	}
	if _, ok := storage[user.Name]; ok {
		return errDuplicate
	}
	storage[user.Name] = user
	return nil
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, errDuplicate):
		http.Error(w, "duplicate user", http.StatusForbidden)
	default:
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

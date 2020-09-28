package test

import (
	"encoding/json"
	"os"
	"testing"
)

// ReadJSON decodes JSON from a file into a value. If it fails to open the given
// filename or decode the contents of the file, ReadJSON will call t.Fatal with the
// error.
//
// The type of the first parameter as testing.TB allows it to be used from both
// tests and benchmarks, passing either a *testing.T or *testing.B.
//
//   cfg := Config{}
//   test.ReadJSON(t, "testdata/config.json", &cfg)
func ReadJSON(t testing.TB, fname string, v interface{}) {
	t.Helper()
	f, err := os.Open(fname) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck,gosec
	if err := json.NewDecoder(f).Decode(v); err != nil {
		t.Fatalf("cannot decode %s: %v", fname, err)
	}
}

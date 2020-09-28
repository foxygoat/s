package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockT struct {
	testing.TB
	errMsg string
}

func (t *MockT) Helper() {}
func (t *MockT) Fatal(s ...interface{}) {
	t.errMsg = fmt.Sprintln(s...)
}

func (t *MockT) Fatalf(format string, args ...interface{}) {
	t.errMsg = fmt.Sprintf(format, args...)
}

func TestReadJSON(t *testing.T) {
	mt := &MockT{TB: t}
	got := map[string]string{}
	ReadJSON(mt, "testdata/profile.json", &got)
	want := map[string]string{"name": "Lee Sedol", "game": "Go"}
	require.Equal(t, want, got)
	require.Empty(t, mt.errMsg)
}

func TestReadJSONMissingFile(t *testing.T) {
	mt := &MockT{}
	got := map[string]string{}
	ReadJSON(mt, "testdata/MISSING_FILE.json", &got)
	require.NotEmpty(t, mt.errMsg)
}

func TestReadJSONInvalidJSON(t *testing.T) {
	mt := &MockT{}
	got := map[string]string{}
	ReadJSON(mt, "testdata/invalid.json", &got)
	require.NotEmpty(t, mt.errMsg)
}

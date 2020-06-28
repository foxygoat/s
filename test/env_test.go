package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetExisting(t *testing.T) {
	err := os.Setenv("FOO", "BAR")
	require.NoError(t, err)

	e := Environ{}
	e2 := e.Set("FOO", "BAZ")
	require.Equal(t, e, e2)

	val, ok := os.LookupEnv("FOO")
	require.True(t, ok, "env var not found in environment: FOO")
	require.Equal(t, "BAZ", val)

	e.Restore()
	val, ok = os.LookupEnv("FOO")
	require.True(t, ok, "env var not found in environment: FOO")
	require.Equal(t, "BAR", val)
}

func TestSetNew(t *testing.T) {
	err := os.Unsetenv("FOO")
	require.NoError(t, err)

	e := Environ{}
	e2 := e.Set("FOO", "BAZ")
	require.Equal(t, e, e2)

	val, ok := os.LookupEnv("FOO")
	require.True(t, ok, "env var not found in environment: FOO")
	require.Equal(t, "BAZ", val)

	e.Restore()
	_, ok = os.LookupEnv("FOO")
	require.False(t, ok, "env var found in environment: FOO")
}

func TestSetPanic(t *testing.T) {
	e := Environ{}
	require.Panics(t, func() { e.Set("FOO=", "BAR") })
	require.Panics(t, func() { e.Set("FOO\x00", "BAR") })
	require.Panics(t, func() { e.Set("FOO", "BAR\x00") })
}

func TestUnsetExisting(t *testing.T) {
	err := os.Setenv("FOO", "BAR")
	require.NoError(t, err)

	e := Environ{}
	e2 := e.Unset("FOO")
	require.Equal(t, e, e2)

	_, ok := os.LookupEnv("FOO")
	require.False(t, ok, "env var found in environment: FOO")

	e.Restore()
	val, ok := os.LookupEnv("FOO")
	require.True(t, ok, "env var not found in environment: FOO")
	require.Equal(t, "BAR", val)
}

func TestUnsetMissing(t *testing.T) {
	err := os.Unsetenv("FOO")
	require.NoError(t, err)

	e := Environ{}
	e2 := e.Unset("FOO")
	require.Equal(t, e, e2)

	_, ok := os.LookupEnv("FOO")
	require.False(t, ok, "env var found in environment: FOO")

	e.Restore()
	_, ok = os.LookupEnv("FOO")
	require.False(t, ok, "env var found in environment: FOO")
}

func TestGlobalEnv(t *testing.T) {
	err := os.Unsetenv("FOO")
	require.NoError(t, err)

	Env.Set("FOO", "BAR")
	val, ok := os.LookupEnv("FOO")
	require.True(t, ok, "env var not found in environment: FOO")
	require.Equal(t, "BAR", val)

	Env.Restore()
	_, ok = os.LookupEnv("FOO")
	require.False(t, ok, "env var found in environment: FOO")
}

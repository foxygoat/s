package test

import (
	"fmt"
	"os"
)

// Environ is a frontend to the OS environment that keeps track of the previous
// value of an environment variable before setting or unsetting it. Later,
// Restore() can be called to reset all modified variables back to their
// original values and their presence.
//
// It is intended to be used in test cases where a particular starting
// environment is needed and to ensure that it is reset back to the starting
// state at the end.
//
// The operations on Environ are not concurrency-safe as the global OS
// environment is also not concurrency-safe. You should not use Environ in test
// functions that are marked t.Parallel().
//
// The methods on Environ return the receiver so environment setup can be
// chained:
//
//     e := test.Environ{}.Set("FOO", "BAR").Unset("BAZ")
//     defer e.Restore()
type Environ map[string]*string

// Env is a global Environ variable for simpler use. Since the OS environment
// is global, using a global test.Environ is rarely a problem. Use is as
// simple as:
//
//     func TestFoo(t *testing.T) {
//         test.Env.Set("FOO", "BAR").Unset("BAZ")
//         defer test.Env.Restore()
//         ... test, test, test ...
//     }
var Env Environ

// Set sets key to value in the OS environment, saving the previous value and
// presence. It returns the receiver so calls can easily be chained.
// If key or value are invalid (as per os.Setenv), this method will panic.
func (e Environ) Set(key, value string) Environ {
	e.save(key)
	set(key, value)
	return e
}

// Unset removes key from the OS environment, saving the previous value and
// presence. It returns the receive so calls can easily be chained. If
// os.Unsetenv returns an error, this method will panic. However, it appears
// that os.Unsetenv does not return an error.
func (e Environ) Unset(key string) Environ {
	e.save(key)
	unset(key)
	return e
}

// Restore restores all the saved values from Set and Unset back to the
// original values and presence. When Restore returns, all saved values will be
// forgotten. This permits the Environ to be reused. If Restore is unable to
// set or unset any variables back to their original state, it will panic.
func (e Environ) Restore() {
	for key, value := range e {
		if value == nil {
			unset(key)
		} else {
			set(key, *value)
		}
		delete(e, key)
	}
}

func set(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		panic(fmt.Errorf("could not Setenv(%#v, %#v): %v", key, value, err))
	}
}

func unset(key string) {
	_ = os.Unsetenv(key)
	// os.Unsetenv never actually returns an error. Comment this out for coverage.
	/*
		if err != nil {
			panic(fmt.Errorf("could not Unsetenv(%#v): %v", key, err))
		}
	*/
}

func (e Environ) save(key string) {
	var old *string
	if v, ok := os.LookupEnv(key); ok {
		old = &v
	}
	e[key] = old
}

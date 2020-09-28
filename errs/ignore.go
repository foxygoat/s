package errs

// Ignore executes nullary input function f and ignores its error return
// value. A typical use case is ignoring errors in the defer statements:
//
//   defer errs.Ignore(body.Close)
func Ignore(f func() error) {
	_ = f()
}

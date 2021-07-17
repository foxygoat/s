package errs

// Const is an Error string that can be used to create a const sentinel
// error like `const ErrFoo = errs.Const("foo")`. As a const, the value
// of the sentinel error cannot be changed.
type Const string

// Error returns c as a string, implementing the error interface.
func (c Const) Error() string {
	return string(c)
}

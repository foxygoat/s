package errs_test

import (
	"errors"
	"fmt"
	"io"

	"foxygo.at/s/errs"
)

func ExampleErrorf() {
	errInternal := errors.New("internal error")
	err := errs.Errorf("%v, caused by: %v", errInternal, io.ErrUnexpectedEOF)

	fmt.Println(err)
	fmt.Println("err is errInternal:", errors.Is(err, errInternal))
	fmt.Println("err is ErrUnexpectedEOF:", errors.Is(err, io.ErrUnexpectedEOF))
	// output: internal error, caused by: unexpected EOF
	// err is errInternal: true
	// err is ErrUnexpectedEOF: true
}

func ExampleNew() {
	errInternal := errors.New("internal error")
	err := errs.New(errInternal, io.ErrUnexpectedEOF)

	fmt.Println(err)
	fmt.Println("err is errInternal:", errors.Is(err, errInternal))
	fmt.Println("err is ErrUnexpectedEOF:", errors.Is(err, io.ErrUnexpectedEOF))
	// output: internal error: unexpected EOF
	// err is errInternal: true
	// err is ErrUnexpectedEOF: true
}

func ExampleNoWrap() {
	errInternal := errors.New("internal error")
	err := errs.Errorf("%v, caused by: %v", errInternal, errs.NoWrap(io.ErrUnexpectedEOF))

	fmt.Println(err)
	fmt.Println("err is errInternal:", errors.Is(err, errInternal))
	fmt.Println("err is ErrUnexpectedEOF:", errors.Is(err, io.ErrUnexpectedEOF))
	// output: internal error, caused by: unexpected EOF
	// err is errInternal: true
	// err is ErrUnexpectedEOF: false
}

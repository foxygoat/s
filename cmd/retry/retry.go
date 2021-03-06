// Command retry re-runs a command a number of times until it succeeds.
//
//   Usage: retry <count> <command> [args...]
//
// retry will run <command> with [args...]. If it exits with a non-zero code,
// it will be re-run up to <count> times. <count> may not be negative.
//
// The exit code of retry is the exit code of the last execution of <command>
// or 1 if retry exits with its own error.
//
// <command> is executed directly with [args...] as provided. If <command> does
// not contain any path separators, the search path is used to locate it. No
// shell is used to run <command>.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	code, err := retry(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}

func retry(args []string) (int, error) {
	if len(args) < 2 {
		return 1, errors.New("usage: retry <count> <command> [args...]")
	}
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return 1, fmt.Errorf("invalid count: %v", err)
	}
	return Retry(i, args[1:]...)
}

// Retry runs a command up to count+1 times (once plus count retries). If the
// command could not be run or it exits with a non-zero error code, then the
// command is retried up to count times. The command's exit code and an error
// (if any) is returned.
func Retry(count int, argv ...string) (int, error) {
	err := errors.New("negative retries are not allowed")
	var cmd *exec.Cmd
	for i := 0; i <= count && err != nil; i++ {
		cmd = exec.Command(argv[0], argv[1:]...) //nolint:gosec
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err = cmd.Run()
	}
	var e *exec.ExitError
	if err != nil && !errors.As(err, &e) {
		return 1, err
	}
	return cmd.ProcessState.ExitCode(), nil
}

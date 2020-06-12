package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	code, err := timeout(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}

func timeout(args []string) (int, error) {
	if len(args) < 2 {
		return 1, errors.New("usage: timeout <duration> command [args...]")
	}
	d, err := time.ParseDuration(args[0])
	if err != nil {
		return 1, fmt.Errorf("invalid duration: %v", err)
	}
	return Timeout(d, args[1:]...)
}

// Timeout runs a command, killing it after a given duration. The command's
// exit code and an error (if any) are returned. If the command could not be
// run, or the timeout expires, the exit code returned is 1.
func Timeout(d time.Duration, argv ...string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...) //nolint:gosec
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return 1, fmt.Errorf("timeout (%s)", d)
	}
	var e *exec.ExitError
	if err != nil && !errors.As(err, &e) {
		return 1, err
	}
	return cmd.ProcessState.ExitCode(), nil
}

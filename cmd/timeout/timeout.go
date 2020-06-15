package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var (
	grace = flag.Duration("g", 3*time.Second, "Grace period before sending SIGKILL")
	sig   = flag.Int("s", int(syscall.SIGTERM), "Signal to send on timeout")
)

func main() {
	flag.Parse()
	code, err := timeout(flag.Args())
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

// Timeout runs a command, signalling it after a given duration. If the process
// has not exited by the time the grace period passes, the process is killed.
// The command's exit code and an error (if any) are returned. If the command
// could not be run, or the timeout expires, the exit code returned is 1.
func Timeout(d time.Duration, argv ...string) (int, error) {
	cmd, err := exec.LookPath(argv[0])
	if err != nil {
		return 1, err
	}
	// Prevent stdin, stdout and stderr from being closed by StartProcess
	attr := &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	p, err := os.StartProcess(cmd, argv, attr)
	if err != nil {
		return 1, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	isTimeout := make(chan bool)
	go func() {
		timedOut := callAfter(ctx, d, func() { _ = p.Signal(syscall.Signal(*sig)) })
		if timedOut {
			callAfter(ctx, *grace, func() { _ = p.Kill() })
		}
		isTimeout <- timedOut
	}()
	ps, err := p.Wait()
	cancel()
	if <-isTimeout {
		return 1, fmt.Errorf("timeout (%s)", d)
	}
	if err != nil {
		return 1, err
	}
	return ps.ExitCode(), nil
}

// callAfter calls function f after duration d and returns true. If the context
// is cancelled before duration d, return false without invoking f.
func callAfter(ctx context.Context, d time.Duration, f func()) bool {
	select {
	case <-time.After(d):
		f()
		return true
	case <-ctx.Done():
		return false
	}
}

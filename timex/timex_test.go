package timex

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestTicker(t *testing.T) {
	ticker := NewTicker(20 * time.Millisecond)
	time.AfterFunc(50*time.Millisecond, ticker.Stop)
	cnt := 0
	for range ticker.C {
		cnt++
	}
	// May be flaky due to OS scheduling
	require.Equal(t, 2, cnt)
}

func TestTickerPanic(t *testing.T) {
	require.Panics(t, func() { NewTicker(-2 * time.Millisecond) })
}

func TestTickerWithContext(t *testing.T) {
	errFive := fmt.Errorf("five")
	cnt := 0
	f := func() error {
		cnt++
		if cnt == 5 {
			return errFive
		}
		return nil
	}

	// errgroup.Group will cancel the context when one of the functions it
	// calls returns an error
	g, ctx := errgroup.WithContext(context.Background())
	ticker := NewTickerWithContext(ctx, 2*time.Millisecond)

	for range ticker.C {
		g.Go(f)
	}
	err := g.Wait()
	require.Error(t, err)
	require.Equal(t, errFive, err)
	require.Equal(t, 5, cnt)
}

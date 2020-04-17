package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTicker(t *testing.T) {
	ticker := NewTicker(2 * time.Millisecond)
	time.AfterFunc(5*time.Millisecond, ticker.Stop)
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

// Package timex provides stdlib time package related utilities.
//
// It exports type Ticker which works like stdlib time.Ticker but closes
// its channel on stop.
package timex

import (
	"context"
	"time"
)

// A Ticker holds a channel that delivers 'ticks' of a clock at
// intervals, similar to time.Ticker from the standard library.
type Ticker struct {
	C    <-chan time.Time
	done chan struct{}
}

// NewTicker returns a new Ticker containing a channel that will send
// the time with a period specified by the duration argument, similar to
// standard library's time.Ticker.
//
// NewTicker will panic if the duration is negative.
//
// It differs from the standard library Ticker by closing C when
// Stop() is called, which allows for ranging over C:
//
//		for range ticker.C {
//   		// do something
//		}
func NewTicker(d time.Duration) *Ticker {
	c := make(chan time.Time)
	done := make(chan struct{})
	tt := time.NewTicker(d)
	go func() {
		for {
			select {
			case tick := <-tt.C:
				c <- tick
			case <-done:
				tt.Stop()
				close(c)
				return
			}
		}
	}()
	return &Ticker{C: c, done: done}
}

// Stop turns off the ticker, closing the ticker channel C, after which
// no more ticks will be sent. It is safe to call Stop more than once;
// subsequent calls to Stop do nothing.
//
// When working with select avoid erroneous ticks on close with:
// 		for {
// 		    select {
// 		    case _, ok := <-ticker.C:
// 		        if ok {
// 		            // do something
// 		        } else {
// 		            return
// 		        }
// 		    }
// 		}
func (t *Ticker) Stop() {
	select {
	case <-t.done:
	default:
		close(t.done)
	}
}

// NewTickerWithContext creates a Ticker that stops when the context is
// cancelled.
func NewTickerWithContext(ctx context.Context, d time.Duration) *Ticker {
	t := NewTicker(d)
	if ctx.Done() != nil {
		go func() {
			<-ctx.Done()
			t.Stop()
		}()
	}
	return t
}

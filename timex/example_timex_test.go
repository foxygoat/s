package timex_test

import (
	"context"
	"fmt"
	"time"

	"foxygo.at/s/timex"
	"golang.org/x/sync/errgroup"
)

func ExampleNewTicker() {
	ticker := timex.NewTicker(20 * time.Millisecond)
	time.AfterFunc(50*time.Millisecond, ticker.Stop)
	cnt := 0
	for range ticker.C {
		cnt++
	}
	fmt.Println("cnt:", cnt)
	// output: cnt: 2
}

func ExampleNewTickerWithContext() {
	cnt := 0
	f := func() error {
		cnt++
		if cnt == 5 {
			return fmt.Errorf("five")
		}
		return nil
	}

	// errgroup.Group will cancel the context when one of the functions it
	// calls returns an error
	g, ctx := errgroup.WithContext(context.Background())
	ticker := timex.NewTickerWithContext(ctx, 2*time.Millisecond)

	for range ticker.C {
		g.Go(f)
	}
	err := g.Wait()
	fmt.Println("cnt:", cnt, "err:", err)
	// output: cnt: 5 err: five
}

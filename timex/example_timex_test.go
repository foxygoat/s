package timex_test

import (
	"fmt"
	"time"

	"foxygo.at/s/timex"
)

func ExampleNewTicker() {
	ticker := timex.NewTicker(2 * time.Millisecond)
	time.AfterFunc(5*time.Millisecond, ticker.Stop)
	cnt := 0
	for range ticker.C {
		cnt++
	}
	fmt.Println("cnt:", cnt)
	// output: cnt: 2
}

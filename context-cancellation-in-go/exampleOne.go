package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func exampleOne() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	wg := new(sync.WaitGroup)
	wg.Add(2) // we will start two ticker goroutines.

	go func() {
		// start the two ticker functions in separate goroutines.
		go runTicker(ctx, wg)
		go anotherTicker(ctx, wg)

		// sleep for 10s for ticker functions to print something to the console.
		time.Sleep(10 * time.Second)

		// after 10s, cancel the context.
		cancelFunc()
	}()

	// wait for the ticker functions to complete.
	wg.Wait()
}

func anotherTicker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-time.After(2 * time.Second):
			sleepContext(ctx, time.Minute) // sleep here indicates some work.
			fmt.Println(time.Now(), "another ticker executed")
		case <-ctx.Done():
			return
		}
	}

}

func runTicker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-time.After(2 * time.Second):
			fmt.Println(time.Now(), "ticker executed")
		case <-ctx.Done():
			return
		}
	}
}

// a sleep function that honors context cancellation.
func sleepContext(ctx context.Context, delay time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(delay):
	}
}

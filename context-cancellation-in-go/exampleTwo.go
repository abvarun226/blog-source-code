package main

import (
	"context"
	"fmt"
	"time"
)

func exampleTwo() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc() // cancel the context immediately.

	for i := 0; i < 5; i++ {
		execute(ctx) // execute function 5 times, each time there is no sleep.
	}
}

func execute(ctx context.Context) {
	fmt.Println(time.Now(), "executing func")

	// since ctx was cancelled immediately, sleepContext will return immediately without any sleep.
	sleepContext(ctx, 5*time.Second)
}

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Example of usage 'context.WithTimeout()' and '<context>.Deadline()'.
func run_context_deadline() {

	wg := &sync.WaitGroup{}

	discounts := func(duration time.Duration) {
		defer wg.Done()

		select {
		case <-time.After(duration):
		}
	}

	shopping := func(ctx context.Context, roadToShopDuration time.Duration) {
		defer wg.Done()

		if deadline, ok := ctx.Deadline(); ok {
			if deadline.Sub(time.Now().Add(roadToShopDuration)) <= 0 {
				fmt.Println("No chance to buy with sale(")
				return
			}
		}

		fmt.Println("Go shopping!")
	}

	discountsDuration := 1 * time.Second
	wg.Add(1)
	go discounts(discountsDuration)

	ctx, _ := context.WithTimeout(context.Background(), discountsDuration)
	roadToShopDuration := 5 * time.Second
	wg.Add(1)
	go shopping(ctx, roadToShopDuration)

	wg.Wait()
}

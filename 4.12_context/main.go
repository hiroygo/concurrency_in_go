package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func workSlowly(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(20 * time.Second):
		fmt.Println("workSlowly done")
		return nil
	}
}

func workQuickly(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Second):
		fmt.Println("workQuickly done")
		return nil
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := workSlowly(ctx); err != nil {
			fmt.Printf("workSlowly error %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(ctx, 5*time.Second)
		if err := workQuickly(ctx); err != nil {
			fmt.Printf("workQuickly error %v\n", err)
			cancel()
		}
	}()

	wg.Wait()
}

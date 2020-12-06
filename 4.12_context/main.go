package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func work(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		fmt.Printf("work %v done\n")
		return nil
	}
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	// キャンセルが実行されてこのゴルーチンもキャンセルされる
	go func() {
		defer wg.Done()
		if err := work(ctx, 10*time.Second); err != nil {
			fmt.Printf("goroutine1: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	// このゴルーチンが先にタイムアウトして
	// キャンセルを実行する
	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(ctx, 5*time.Second)
		if err := work(ctx, 15*time.Second); err != nil {
			fmt.Printf("goroutine2: %v\n", err)
			cancel()
		}
	}()

	wg.Wait()
}

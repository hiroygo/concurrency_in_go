package main

import (
	"fmt"
	"math/rand"
)

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for {
			select {
			case <-done:
				return
			case c <- fn():
			}
		}
	}()
	return c
}

func take(done <-chan interface{}, in <-chan interface{}, n int) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				c <- v
			}
		}
	}()
	return c
}

func main() {
	done := make(chan interface{})
	defer close(done)

	randCh := take(done, repeatFn(done, func() interface{} { return rand.Int() }), 10)
	for v := range randCh {
		fmt.Println(v)
	}
}

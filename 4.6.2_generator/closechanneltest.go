package main

import (
	"fmt"
)

func generator(done <-chan interface{}, args ...interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for _, a := range args {
			select {
			case <-done:
				return
			case c <- a:
			}
		}
	}()
	return c
}

// テスト用に done を削除する
func nooper(in <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)

		// 以下だと in が閉じられても c に延々と nil が送り込まれる
		// for {
		//  c <- <-in
		// }
		for {
			select {
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

	ch := nooper(generator(done, 1, 2, 3, 4, 5))
	for v := range ch {
		fmt.Println(v)
	}
}

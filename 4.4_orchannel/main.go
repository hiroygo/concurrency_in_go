package main

import (
	"fmt"
	"time"
)

func main() {
	// 34 行目で再帰呼び出しするために初めに宣言が必要
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		// case 0 は再帰では使わないはず
        // done がいるので
		// 多分ガード用
		case 0:
			fmt.Printf("return 0 %v\n", len(channels))
			return nil
		case 1:
			fmt.Printf("return 1 %v\n", len(channels))
			return channels[0]
		}

		done := make(chan interface{})
		go func() {
			defer close(done)
			defer fmt.Printf("bye %v\n", len(channels))

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				// slice[3:] では len(slice) < 4 であっても空のスライスが返る
				case <-or(append(channels[3:], done)...):
				}
			}
		}()
		return done
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	<-or(
		sig(1*time.Second),
		sig(5*time.Minute),
	)
	//	<-or(
	//		sig(2*time.Hour),
	//		sig(5*time.Minute),
	//		sig(1*time.Second),
	//	)
	time.Sleep(5 * time.Second)
}

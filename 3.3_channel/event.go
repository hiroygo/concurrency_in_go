package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan interface{})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// ゴルーチン内で i := i のシャドーイングは使えない
			// ゴルーチンの起動時には i が違う値になっている可能性があるため
			// ゴルーチンではなく、ただの関数実行なら問題ない
			<-ch
			fmt.Println(i)
		}(i)
	}

	for i := 0; i < 3; i++ {
		fmt.Println("sleeping...")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("go!")
	close(ch)
	wg.Wait()
}

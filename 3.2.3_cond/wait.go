package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Mutex の初期状態は unlocked されている
	c := sync.NewCond(&sync.Mutex{})

	for i := 0; i < 5; i++ {
		go func() {
			time.Sleep(1 * time.Second)

			// 以下は不要なら実行しなくていい
			// c.L.Lock()
			// c.L.Unlock()
			i := i
			fmt.Printf("go %v\n", i)
			c.Signal()
		}()

		// Wait を呼び出すと内部で Unlock が実行される
		// Unlock されている Mutex を Unlock するとエラーになる
		c.L.Lock()
		c.Wait()
		// Wait は終了時に内部で Lock を呼び出す
		c.L.Unlock()
	}

	fmt.Println("bye")
}

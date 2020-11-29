package main

import "fmt"

func main() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var n1, n2 int
	for i := 0; i < 1000; i++ {
		// go のランタイムは case の選択に疑似乱数による
		// 一様選択をしている
		select {
		case <-c1:
			n1++
		case <-c2:
			n2++
		}
	}

	// e.g. `n1:477, n2:523`
	fmt.Printf("n1:%d, n2:%d\n", n1, n2)
}

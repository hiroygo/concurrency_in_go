package main

import (
	"fmt"
)

// チャネルの make から close まで行う関数
func wordGenarator() <-chan string {
	words := []string{"dog", "panda", "orange", "book", "curry"}
	// あらかじめ送信する要素数がわかっているので、バッファを確保しておく
	// 書籍だとなぜか `len - 1` していた
	ch := make(chan string, len(words))

	go func() {
		// 確実に閉じる
		defer close(ch)
		defer fmt.Println("goroutine end")

		for i := 0; i < len(words); i++ {
			ch <- words[i]
		}
	}()

	return ch
}

func main() {
	// チャネルの読み込みだけ行う
	// チャネルの make,close と読み込みだけ行う部分を分けて責任を分割する
	wordGen := wordGenarator()
	// 以下はエラーになる。読み込み専用チャネルを閉じることはできない
	// close(wordGen)

	for w := range wordGen {
		fmt.Println(w)
	}
	fmt.Println("main end")
}

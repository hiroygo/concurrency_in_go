package main

import "fmt"

func orDone(done <-chan interface{}, in <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for {
			select {
			case <-done:
				return
			case s, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case c <- s:
				}
			}
		}
	}()
	return c
}

func generator(done <-chan interface{}, strs ...string) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for _, s := range strs {
			select {
			case <-done:
				return
			case c <- s:
			}
		}
	}()
	return c
}

// rune チャネルを送るチャネルを返す
func chanGenerator(done <-chan interface{}, in <-chan string) <-chan <-chan rune {
	cc := make(chan (<-chan rune))
	go func() {
		defer close(cc)
		for s := range orDone(done, in) {
			c := make(chan rune)
			go func(str string) {
				defer close(c)
				for _, r := range str {
					select {
					case <-done:
						return
					case c <- r:
					}
				}
			}(s)

			select {
			case <-done:
				return
			case cc <- c:
			}
		}
	}()
	return cc
}

// 複数の rune チャネルを 1 つの rune チャネルにまとめる
func bridge(done <-chan interface{}, cc <-chan <-chan rune) <-chan rune {
	c := make(chan rune)
	go func() {
		defer close(c)
		for {
			select {
			case <-done:
				return
			case runeCh, ok := <-cc:
				if !ok {
					return
				}

			exit_loop:
				for {
					select {
					case <-done:
						return
					case r, ok := <-runeCh:
						if !ok {
							break exit_loop
						}

						select {
						case <-done:
							return
						case c <- r:
						}
					}
				}
			}
		}
	}()
	return c
}

func main() {
	done := make(chan interface{})
	defer close(done)

	gen := generator(done, "dog", "orange")
	chans := chanGenerator(done, gen)
	runeCh := bridge(done, chans)
	for r := range runeCh {
		fmt.Printf("%c", r)
	}
	fmt.Println()
}

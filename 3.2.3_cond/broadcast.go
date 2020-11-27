package main

import (
	"fmt"
	"sync"
)

type button struct {
	clicked *sync.Cond
}

func (b *button) subscribe(fn func()) {
	var goroutineRunning sync.WaitGroup
	goroutineRunning.Add(1)
	go func() {
		goroutineRunning.Done()
		b.clicked.L.Lock()
		defer b.clicked.L.Unlock()
		b.clicked.Wait()
		fn()
	}()
	goroutineRunning.Wait()
}

func (b *button) Broadcast() {
	b.clicked.Broadcast()
}

func main() {
	btn := button{clicked: sync.NewCond(&sync.Mutex{})}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	btn.subscribe(func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})
	btn.subscribe(func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})
	btn.subscribe(func() {
		fmt.Println("Displaying annoying dialogue box!")
		clickRegistered.Done()
	})

	btn.Broadcast()
	clickRegistered.Wait()
}
